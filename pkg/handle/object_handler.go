package handle

import (
	"assessment/pkg/get"
	"assessment/pkg/persist"
	"time"
)

type ObjectHandler struct {
	persistence   persist.Persistor
	objectGetter  get.Getter
	durationValid time.Duration
}

func NewObjectHandler(pers persist.Persistor, getter get.Getter, durationValid time.Duration) *ObjectHandler {
	return &ObjectHandler{
		durationValid: durationValid,
		persistence:   pers,
		objectGetter:  getter,
	}
}
