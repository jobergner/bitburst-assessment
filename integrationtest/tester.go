package main

import (
	"assessment/pkg/persist"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type integrationTester struct {
	tp testPoster
	pg persist.Persistor
	t  *testing.T
}

func (it *integrationTester) postAndExpect(idsToPost, idsToExpect []int) {
	it.tp.postIDs(idsToPost...)
	time.Sleep(time.Millisecond * 5)
	currentObjects, err := it.pg.GetObjects()
	assert.Nil(it.t, err)
	assert.ElementsMatch(it.t, idsToExpect, listObjectIDs(currentObjects))
}

func newTester(t *testing.T) *integrationTester {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(fmt.Errorf("error reading env file: %s", err))
	}
	pg := persist.NewPostgres()
	if err := pg.Connect(); err != nil {
		t.Fatalf("could not connect to postgres: %s", err)
	}
	client := &http.Client{Timeout: 1 * time.Second}
	tp := testPoster{client}
	return &integrationTester{
		pg: pg,
		tp: tp,
		t:  t,
	}
}
