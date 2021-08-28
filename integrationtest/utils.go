package integrationtest

import (
	"assessment/pkg/object"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	servicePort      = ":9090"
	objectServerPort = ":9010"
)

func serveObjectSource() *http.Server {
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
		Addr:         objectServerPort,
		Handler:      mux,
	}

	go func() { _ = srv.ListenAndServe() }()

	return srv
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
	resp, err := tp.c.Post(fmt.Sprintf("http://localhost%s/callback", servicePort), "application/json", body)
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

func waitForIt(url string) bool {

	for i := 0; i < 5; i++ {

		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 404 {
				return true
			}
		}

		time.Sleep(time.Second / 2)
	}

	return false
}
