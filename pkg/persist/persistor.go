package persist

import "assessment/pkg/object"

type Persistor interface {
	WriteObject(object.Object) error
	GetObjects() (map[int]object.Object, error)
	DeleteObject(int, int64) error
	Connect() error
}
