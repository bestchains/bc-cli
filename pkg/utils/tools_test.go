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

package utils

import (
	"reflect"
	"testing"
)

func TestGetNestedString(t *testing.T) {
	obj := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}

	// Test a valid nested string value.
	expected := "baz"
	actual := GetNestedString(obj, "foo", "bar")
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	// Test a missing nested field.
	expected = ""
	actual = GetNestedString(obj, "foo", "notfound")
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	// Test a non-string nested value.
	obj["foo"].(map[string]interface{})["baz"] = 123
	expected = ""
	actual = GetNestedString(obj, "foo", "baz")
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestRemoveDuplicateForStringSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "some duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a", "a"},
			expected: []string{"a"},
		},
		{
			name:     "empty string",
			input:    []string{"a", "b", "", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RemoveDuplicateForStringSlice(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
