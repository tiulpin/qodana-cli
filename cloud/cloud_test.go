/*
 * Copyright 2021-2023 JetBrains s.r.o.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cloud

import (
	"net/http"
	"testing"
)

func TestGetProjectByBadToken(t *testing.T) {
	t.Skip() // Until qodana.cloud response is fixed
	client := NewQdClient("https://www.jetbrains.com")
	result := client.getProject()
	switch v := result.(type) {
	case Success:
		t.Errorf("Did not expect request error: %v", v)
	case APIError:
		if v.StatusCode > http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, v.StatusCode)
		}
	case RequestError:
		t.Errorf("Did not expect request error: %v", v)
	default:
		t.Error("Unknown result type")
	}
}

func TestValidateToken(t *testing.T) {
	client := NewQdClient("kek")
	if projectName := client.ValidateToken(); projectName != "" {
		t.Errorf("Problem")
	}
}
