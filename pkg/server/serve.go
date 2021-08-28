package server

import (
	"fmt"
	"net/http"
	"time"
)

type ServerConfig struct {
	UseMockDB      bool
	ObjectLifespan int64
	ObjectSource   string
}

func Serve(addr string, c ServerConfig) *http.Server {

	h := configureObjectHandler(c.UseMockDB, c.ObjectSource, c.ObjectLifespan)

	mux := http.NewServeMux()
	mux.Handle("/callback", callbackHandler(h))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
		Handler:      mux,
	}

	fmt.Println("server running at", addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	return srv
}
