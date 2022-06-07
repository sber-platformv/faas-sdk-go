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
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	writeTimeout    = time.Second * 60
	readTimeout     = time.Second * 60
	gracefulTimeout = time.Second * 15
)

func Start(port string, function func(http.ResponseWriter, *http.Request)) {
	router := http.NewServeMux()
	router.HandleFunc("/", function)
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		Handler:        router,
	}

	go func() {
		log.Print("Serving function")
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Print("HTTP server closed")
		} else {
			log.Fatalf("can't start listen: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	log.Printf("Received signal %s - shutting down...", sig.String())

	_ = srv.Shutdown(ctx)

	log.Print("shutting down")
}