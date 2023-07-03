/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package depository

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/bestchains/bc-cli/pkg/auth"
	"github.com/bestchains/bc-cli/pkg/common"
	"github.com/bestchains/bc-cli/pkg/utils"
)

func download(host, style string, kids []string, option common.Options) {

	limit := make(chan struct{}, 3)
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))

	for _, kid := range utils.RemoveDuplicateForStringSlice(kids) {
		wg.Add(1)
		go func(kid string) {
			defer wg.Done()
			limit <- struct{}{}
			defer func() {
				<-limit
			}()

			name := kid + ".pdf"
			u := fmt.Sprintf("%s%s?style=%s", host, fmt.Sprintf(common.DepositoryCertificate, kid), style)
			req, err := http.NewRequest(http.MethodGet, u, nil)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "new request for %s error %s", name, err)
				return
			}
			auth.AddAuthHeader(req)
			client := &http.Client{
				// the size mf the certificate file is usually around 3-10M, so the 10s timeout time is more reasonable.
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "do request for %s error %s", name, err)
				return
			}
			defer resp.Body.Close()

			buf := make([]byte, 512)
			contentLength := resp.ContentLength
			bar := p.AddBar(contentLength,
				mpb.PrependDecorators(decor.Name(fmt.Sprintf("Downloading %s", name))),
				mpb.BarWidth(50),
				mpb.AppendDecorators(decor.Percentage(decor.WCSyncSpace)))

			f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(option.ErrOut, "open or create file %s error %s", name, err)
				return
			}
			defer f.Close()
			for {
				n, err := resp.Body.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Fprintf(option.ErrOut, "read %s's body error %s", name, err)
						return
					}
				}
				if n > 0 {
					_, _ = f.Write(buf[:n])
					bar.EwmaIncrInt64(int64(n), time.Second)
					continue
				}
				return
			}
		}(kid)
	}
	p.Wait()
}
