package handle

import (
	"assessment/pkg/get"
	"assessment/pkg/persist"
	"time"
)

type Handler struct {
	persistence   persist.Persistor
	objectGetter  get.Getter
	durationValid time.Duration
}

func NewHandler(pers persist.Persistor, getter get.Getter, durationValid time.Duration) *Handler {
	return &Handler{
		durationValid: durationValid,
		persistence:   pers,
		objectGetter:  getter,
	}
}
