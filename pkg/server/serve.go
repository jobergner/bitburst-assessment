package server

import (
	"fmt"
	"net/http"
	"time"
)

type ServerConfig struct {
	UseMockDB      bool
	ObjectLifespan time.Duration
	ObjectSource   string
}

func Serve(addr string, c ServerConfig) *http.Server {

	h := configureObjectHandler(c.UseMockDB, c.ObjectSource, c.ObjectLifespan)

	err := h.ClearExpiredObjects()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/callback", callbackHandler(h))

	srv := &http.Server{
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
		Handler:      mux,
	}

	fmt.Println("server running at", addr)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return srv
}
