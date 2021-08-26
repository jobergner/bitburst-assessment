package main

import (
	"assessment/pkg/get"
	"assessment/pkg/handle"

	"github.com/joho/godotenv"

	"assessment/pkg/persist"
	"flag"
	"fmt"
	"net/http"
	"time"
)

const (
	mockObjectSourceIdentifier string = "OBJECT_SOURCE_MOCK"
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

	var persistence persist.Persistor
	if !*useMockDB {
		persistence = persist.NewPostgres()
	} else {
		persistence = persist.NewMockPersistence()
	}

	err = persistence.Connect()
	if err != nil {
		panic(err)
	}

	var getter get.Getter
	if *objectSource != mockObjectSourceIdentifier {
		getter = get.NewRemoteObjectGetter(*objectSource)
	} else {
		getter = &get.MockObjectGetter{}
	}

	h := handle.NewHandler(persistence, getter, time.Second*time.Duration(*objectLifespan))

	http.HandleFunc("/callback", callbackHandler(h))

	fmt.Println("server running at :9090")
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		panic((fmt.Errorf("error serving: %s", err)))
	}
}
