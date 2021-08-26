package main

import (
	"assessment/pkg/handle"
	"assessment/pkg/object"
	"assessment/pkg/persist"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type callbackBody struct {
	ObjectIDs []int `json:"object_ids"`
}

func callbackHandler(h *handle.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), 500)
			return
		}
		defer r.Body.Close()

		var c callbackBody
		err = json.Unmarshal(b, &c)
		if err != nil {
			http.Error(w, fmt.Sprintf("error unmarshalling received body `%s`: %s", string(b), err.Error()), 500)
			return
		}

		for _, objectID := range c.ObjectIDs {
			go h.Handle(objectID)
		}
	}
}

func main() {
	postgres := persist.NewPostgres()
	err := postgres.Connect()
	if err != nil {
		panic(err)
	}

	err = postgres.WriteObject(object.Object{ObjectID: 2, Online: true, LastSeen: 2, ValidUntil: 3})
	if err != nil {
		panic(err)
	}

	objs, err := postgres.GetObjects()
	if err != nil {
		panic(err)
	}
	fmt.Println(objs)

	err = postgres.DeleteObject(1, 2)
	if err != nil {
		panic(err)
	}

	// mockPers := persist.NewMockPersistence()
	// mockPers.Connect()
	// mockGetter := get.MockObjectGetter{}
	// h := handle.NewHandler(mockPers, &mockGetter, time.Second*5)

	// http.HandleFunc("/callback", callbackHandler(h))

	// err := http.ListenAndServe(":9090", nil)
	// if err != nil {
	// 	panic((fmt.Errorf("error serving: %s", err)))
	// }
}
