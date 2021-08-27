package handle

import (
	"assessment/pkg/object"
	"log"
	"time"
)

func (h *ObjectHandler) Handle(objectID int) {
	receivedAt := time.Now()

	o, err := h.objectGetter.Get(objectID)
	if err != nil {
		log.Println(err)
		return
	}

	if !o.Online {
		return
	}

	o.LastSeen = receivedAt.UnixNano()
	o.ValidUntil = receivedAt.Add(h.durationValid).UnixNano()

	err = h.persistence.WriteObject(o)
	if err != nil {
		log.Println(err)
		return
	}

	go h.deleteWhenExpired(o)
}

func (h ObjectHandler) deleteWhenExpired(o object.Object) {
	remainingLifespan := o.ValidUntil - time.Now().UnixNano()
	time.Sleep(time.Duration(remainingLifespan * int64(time.Nanosecond)))
	h.persistence.DeleteObject(o.ObjectID, o.LastSeen)
}
