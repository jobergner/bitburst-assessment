package handle

import (
	"assessment/pkg/object"
	"time"
)

type persistor interface {
	WriteObject(object.Object) error
	GetObjects() (map[int]object.Object, error)
	DeleteObject(int, int64) error
}

type objectGetter interface {
	Get(int) (object.Object, error)
}

type Handler struct {
	persistence   persistor
	objectGetter  objectGetter
	durationValid time.Duration
}

func NewHandler(pers persistor, getter objectGetter, durationValid time.Duration) *Handler {
	return &Handler{
		durationValid: durationValid,
		persistence:   pers,
		objectGetter:  getter,
	}
}
