package main

import (
	"assessment/pkg/get"
	"assessment/pkg/handle"
	"assessment/pkg/persist"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

	mockPers := persist.NewMockPersistence()
	mockGetter := get.MockObjectGetter{}
	h := handle.NewHandler(mockPers, &mockGetter, time.Second*5)

	http.HandleFunc("/callback", callbackHandler(h))

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic((fmt.Errorf("error serving: %s", err)))
	}
}
