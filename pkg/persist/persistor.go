package persist

import (
	"assessment/pkg/object"
	"time"
)

type Persistor interface {
	WriteObject(object.Object) error
	GetObjects() (map[int]object.Object, error)
	DeleteObject(int, int64) error
	DeleteObjectsOlderThan(time.Duration) error
	Connect() error
}
