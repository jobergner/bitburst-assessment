package handle

import (
	"assessment/pkg/get"
	"assessment/pkg/object"
	"assessment/pkg/persist"

	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	t.Run("deletes objects after the specified time", func(t *testing.T) {
		mockPers := persist.NewMockPersistence()
		mockGetter := get.MockObjectGetter{}
		h := NewHandler(mockPers, &mockGetter, time.Second)

		h.Handle(1)
		assert.ElementsMatch(t, []int{1}, listCurrentObjectIDs(h))

		h.Handle(2)
		assert.ElementsMatch(t, []int{1, 2}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		h.Handle(3)
		assert.ElementsMatch(t, []int{1, 2, 3}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		assert.ElementsMatch(t, []int{3}, listCurrentObjectIDs(h))
	})

	t.Run("resets timer if an object is seen again", func(t *testing.T) {
		mockPers := persist.NewMockPersistence()
		mockGetter := get.MockObjectGetter{}
		h := NewHandler(mockPers, &mockGetter, time.Second)

		h.Handle(1)
		h.Handle(2)
		assert.ElementsMatch(t, []int{1, 2}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		h.Handle(1)

		time.Sleep(time.Second / 2)

		assert.ElementsMatch(t, []int{1}, listCurrentObjectIDs(h))
	})
}

func containsExpectedObjects(objs map[int]object.Object, shouldContain ...int) bool {
	if len(objs) != len(shouldContain) {
		return false
	}

	for expectedID := range shouldContain {
		_, ok := objs[expectedID]
		if !ok {
			return false
		}
	}

	return true
}

func listCurrentObjectIDs(h *Handler) []int {
	var ids []int

	objs, _ := h.persistence.GetObjects()

	for objectID := range objs {
		ids = append(ids, objectID)
	}

	return ids
}
