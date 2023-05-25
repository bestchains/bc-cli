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
	"fmt"
	"strconv"
	"time"
)

type Depository struct {
	Index       string `json:"index" pg:"index"`
	KID         string `json:"kid" pg:"kid,pk"`
	Platform    string `json:"platform" pg:"platform"`
	Operator    string `json:"operator" pg:"operator"`
	Owner       string `json:"owner" pg:"owner"`
	BlockNumber uint64 `json:"blockNumber" pg:"blockNumber"`

	// Content related
	Name             string `json:"name" pg:"name"`
	ContentName      string `json:"contentName" pg:"contentName"`
	ContentID        string `json:"contentID" pg:"contentID"`
	ContentType      string `json:"contentType" pg:"contentType"`
	TrustedTimestamp string `json:"trustedTimestamp" pg:"trustedTimestamp"`
}

func (d Depository) GetByHeader(s string) string {
	switch s {
	case "index":
		return d.Index
	case "kid":
		return d.KID
	case "platform":
		return d.Platform
	case "operator":
		return d.Operator
	case "owner":
		return d.Owner

	case "blockNumber":
		return fmt.Sprintf("%d", d.BlockNumber)
	case "name":
		return d.Name
	case "contenatName":
		return d.ContentName
	case "id", "contentID":
		return d.ContentID
	case "time", "trustedTimestamp":
		seconds, err := strconv.Atoi(d.TrustedTimestamp)
		if err != nil {
			return d.TrustedTimestamp
		}
		now := time.Unix(int64(seconds), 0)
		return now.Format("2006-01-02T15:04:05")
	case "type", "contentType":
		return d.ContentType
	}
	return "<none>"
}

type ValueDepository struct {
	Name             string `json:"name"`
	ContentType      string `json:"contentType"`
	ContentID        string `json:"contentID"` // hash of the file
	TrustedTimestamp string `json:"trustedTimestamp"`
	Platform         string `json:"platform"`
}
