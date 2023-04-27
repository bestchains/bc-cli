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
	"bytes"
	"fmt"
	"testing"
)

type OnlyHeader struct{}

func (OnlyHeader) GetByHeader(string) string {
	return ""
}
func TestOnlyHeader(t *testing.T) {
	headers := []string{"Field1", "Field2", "Field3"}
	expect := "FIELD1    FIELD2    FIELD3\n"
	w := bytes.NewBuffer([]byte{})
	Print(w, headers, []Printer{})
	if got := w.String(); got != expect {
		t.Fatalf("expect '%s' get '%s'", expect, got)
	}
	w1 := bytes.NewBuffer([]byte{})
	expect1 := "FIELD1    FIELD2    FIELD3\n                    \n"
	objs := []Printer{OnlyHeader{}}
	Print(w1, headers, objs)
	if got := w1.String(); got != expect1 {
		t.Fatalf("expect '%s' get '%s'", expect1, got)
	}
}

type OneFiled struct {
	Name string
}

func (one OneFiled) GetByHeader(f string) string {
	return one.Name
}

func TestOneField(t *testing.T) {
	headers := []string{"name"}
	expect := "NAME\nabc\ndef\n"
	w := bytes.NewBuffer([]byte{})
	Print(w, headers, []Printer{OneFiled{Name: "abc"}, OneFiled{Name: "def"}})
	if got := w.String(); got != expect {
		t.Fatalf("expect '%s' get '%s'", expect, got)
	}
}

type ManyFields struct {
	Name  string
	Index int
	X     string
}

func (m ManyFields) GetByHeader(f string) string {
	switch f {
	case "name":
		return m.Name
	case "index":
		return fmt.Sprintf("%d", m.Index)
	case "x":
		return m.X
	}
	return "<none>"
}

func TestManyField(t *testing.T) {
	withoutMissingFiled := []string{"name", "index", "x"}
	withMissingFiled := []string{"name", "index", "x", "y"}

	objs := []Printer{ManyFields{Name: "abc", Index: 1, X: "x1"}, ManyFields{Name: "def", Index: 2, X: "x2"}, ManyFields{Name: "ghi", Index: 3, X: "x3"}}
	expect1 := "NAME    INDEX    X\nabc     1        x1\ndef     2        x2\nghi     3        x3\n"
	w1 := bytes.NewBuffer([]byte{})
	Print(w1, withoutMissingFiled, objs)
	if got := w1.String(); got != expect1 {
		t.Fatalf("expect '%s'  get '%s'", expect1, got)
	}

	expect2 := "NAME    INDEX    X     Y\nabc     1        x1    <none>\ndef     2        x2    <none>\nghi     3        x3    <none>\n"
	w2 := bytes.NewBuffer([]byte{})
	Print(w2, withMissingFiled, objs)
	if got := w2.String(); got != expect2 {
		t.Fatalf("expect '%q' get '%q'", expect2, got)
	}
}
