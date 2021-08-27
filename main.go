package main

import (
	"github.com/joho/godotenv"

	"flag"
	"fmt"
	"net/http"
	"time"
)

const (
	mockObjectSourceIdentifier string = "OBJECT_SOURCE_MOCK"
	addr                       string = ":9090"
)

var objectSource = flag.String("src", mockObjectSourceIdentifier, "endpoint url for object source, uses mock when not specified")
var useMockDB = flag.Bool("mock_db", false, "whether to use the postgres db or a in-memory mock")
var objectLifespan = flag.Int64("ol", 30, "how long an object will be persisted in seconds, defaults to 30")

func main() {
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("error reading env file: %s", err))
	}

	h := configureObjectHandler()

	mux := http.NewServeMux()
	mux.Handle("/callback", callbackHandler(h))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
		Handler:      mux,
	}

	fmt.Println("server running at", addr)
	err = srv.ListenAndServe()
	if err != nil {
		panic((fmt.Errorf("error serving: %s", err)))
	}
}
