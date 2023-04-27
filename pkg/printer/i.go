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

package printer

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Printer print the contents of the structure based on the given header information
type Printer interface {
	GetByHeader(string) string
}

func Print(output io.Writer, headers []string, objs []Printer) {
	w := tabwriter.NewWriter(output, 1, 1, 4, ' ', 0)
	headersCopy := make([]string, len(headers))
	for i := 0; i < len(headers); i++ {
		headersCopy[i] = strings.ToUpper(headers[i])
	}

	fmt.Fprintln(w, strings.Join(headersCopy, "\t"))
	row := make([]string, len(headers))
	for _, o := range objs {
		for i := 0; i < len(headers); i++ {
			row[i] = o.GetByHeader(headers[i])
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}
