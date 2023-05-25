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
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestValueDepotGen(t *testing.T) {
	expectValDepot := ValueDepository{
		Name:             "test name",
		ContentType:      "test type",
		ContentID:        "test ID",
		TrustedTimestamp: "123456789",
		Platform:         "test platform",
	}

	genBase := generateValueDepotBase64("test name", "test type", "test ID", "test platform")

	decodeRes, err := base64.StdEncoding.DecodeString(genBase)
	if err != nil {
		t.Fatalf("decode generated base64 failed: " + err.Error())
	}

	var resValDepot ValueDepository
	err = json.Unmarshal(decodeRes, &resValDepot)
	if err != nil {
		t.Fatalf("unmarshal valDep from decoded base64 failed: " + err.Error())
	}

	if resValDepot.Name != expectValDepot.Name {
		t.Fatalf("generated valDep don't match. expect name '%s', got '%s'", expectValDepot.Name, resValDepot.Name)
	}
}
