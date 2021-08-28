package integrationtest

import (
	"assessment/pkg/persist"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type integrationTester struct {
	tp testPoster
	pg persist.Persistor
	t  *testing.T
}

func (it *integrationTester) postAndExpect(idsToPost, idsToExpect []int) {
	err := it.tp.postIDs(idsToPost...)
	if err != nil {
		it.t.Fatalf(err.Error())
	}
	time.Sleep(time.Millisecond * 100)
	currentObjects, err := it.pg.GetObjects()
	if err != nil {
		it.t.Fatalf(err.Error())
	}
	assert.ElementsMatch(it.t, idsToExpect, listObjectIDs(currentObjects))
}

func newTester(t *testing.T) (*integrationTester, error) {
	pg := persist.NewPostgres()
	if err := pg.Connect(); err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Second}
	tp := testPoster{client}
	return &integrationTester{
		pg: pg,
		tp: tp,
		t:  t,
	}, nil
}
