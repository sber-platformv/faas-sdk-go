// Copyright 2022 АО «СберТех»
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func wrapHTTPFunction(path string, fn func(http.ResponseWriter, *http.Request)) (http.Handler, error) {
	h := http.NewServeMux()
	h.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	})
	return h, nil
}

func TestHTTPFunction(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(w http.ResponseWriter, r *http.Request)
		wantStatus int // defaults to http.StatusOK
		wantResp   string
	}{
		{
			name: "helloworld",
			fn: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "Hello World!")
			},
			wantResp: "Hello World!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, err := wrapHTTPFunction("/", tc.fn)
			if err != nil {
				t.Fatalf("registerHTTPFunction(): %v", err)
			}

			srv := httptest.NewServer(h)
			defer srv.Close()

			resp, err := http.Get(srv.URL)
			if err != nil {
				t.Fatalf("http.Get: %v", err)
			}

			if tc.wantStatus == 0 {
				tc.wantStatus = http.StatusOK
			}
			if resp.StatusCode != tc.wantStatus {
				t.Errorf("TestHTTPFunction status code: got %d, want: %d", resp.StatusCode, tc.wantStatus)
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("ioutil.ReadAll: %v", err)
			}

			if got := strings.TrimSpace(string(body)); got != tc.wantResp {
				t.Errorf("TestHTTPFunction: got %q; want: %q", got, tc.wantResp)
			}
		})
	}
}
