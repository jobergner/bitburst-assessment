package main

import (
	"assessment/pkg/handle"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type callbackBody struct {
	ObjectIDs []int `json:"object_ids"`
}

func callbackHandler(h *handle.ObjectHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), 500)
			return
		}

		var c callbackBody
		err = json.Unmarshal(b, &c)
		if err != nil {
			http.Error(w, fmt.Sprintf("error unmarshalling received body `%s`: %s", string(b), err.Error()), 500)
			return
		}

		for _, objectID := range c.ObjectIDs {
			go h.Handle(objectID)
		}
	})
}
