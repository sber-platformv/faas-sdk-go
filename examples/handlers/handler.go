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

package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Метод Handler. Данный метод будет обрабатывать HTTP запросы поступающие к функции
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	// Чтение тела запроса
	body, _ := ioutil.ReadAll(r.Body)

	// Логирование входящего запроса
	log.Printf("Request received: %s\nMethod: %s", string(body), r.Method)

	// Подготовка и возврат ответа на вызов
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello from Go function!\nYou said: %s", string(body))))
}
