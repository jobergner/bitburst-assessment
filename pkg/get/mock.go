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
		Online:   rand.Float64() > 0.5,
	}, nil
}
