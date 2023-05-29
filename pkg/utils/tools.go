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

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// GetNestedString returns the string value of a nested field.
// Returns "" if value is not found or not a string.
func GetNestedString(obj map[string]interface{}, fields ...string) string {
	val, _, _ := unstructured.NestedString(obj, fields...)
	return val
}

// RemoveDuplicateForStringSlice returns a new slice with duplicate elements removed.
func RemoveDuplicateForStringSlice(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, element := range elements {
		if _, ok := temp[element]; !ok && element != "" {
			temp[element] = struct{}{}
			result = append(result, element)
		}
	}
	return result
}
