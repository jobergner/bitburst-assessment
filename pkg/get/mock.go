package get

import (
	"assessment/pkg/object"
	"math/rand"
	"time"
)

type MockObjectGetter struct {
}

func (r MockObjectGetter) Get(objectID int) (object.Object, error) {
	rand.Seed(time.Now().UnixNano())
	return object.Object{
		ObjectID: objectID,
		Online:   objectID%2 == 0,
	}, nil
}
