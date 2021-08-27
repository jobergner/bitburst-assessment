package main

import (
	"assessment/pkg/object"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func serve() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/objects/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idRaw := strings.TrimPrefix(r.URL.Path, "/objects/")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		w.Write([]byte(fmt.Sprintf(`{"id":%d,"online":%v}`, id, id%2 == 0)))
	}))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":9010",
		Handler:      mux,
	}

	go func() { _ = srv.ListenAndServe() }()

	return srv
}

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	fmt.Println("closing")
}

type testPoster struct {
	c *http.Client
}

func (tp testPoster) postIDs(ids ...int) error {
	var idStrings []string
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"object_ids":[%s]}`, strings.Join(idStrings, ","))))
	resp, err := tp.c.Post("http://localhost:9090/callback", "application/json", body)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func listObjectIDs(objs map[int]object.Object) []int {
	var ids []int

	for objectID := range objs {
		ids = append(ids, objectID)
	}

	return ids
}
