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
		h := NewObjectHandler(mockPers, &mockGetter, time.Second)

		h.Handle(2)
		assert.ElementsMatch(t, []int{2}, listCurrentObjectIDs(h))

		h.Handle(4)
		assert.ElementsMatch(t, []int{2, 4}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		h.Handle(6)
		assert.ElementsMatch(t, []int{2, 4, 6}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		assert.ElementsMatch(t, []int{6}, listCurrentObjectIDs(h))
	})

	t.Run("resets timer if an object is seen again", func(t *testing.T) {
		mockPers := persist.NewMockPersistence()
		mockGetter := get.MockObjectGetter{}
		h := NewObjectHandler(mockPers, &mockGetter, time.Second)

		h.Handle(2)
		h.Handle(4)
		assert.ElementsMatch(t, []int{2, 4}, listCurrentObjectIDs(h))

		time.Sleep(time.Second / 2)

		h.Handle(2)

		time.Sleep(time.Second / 2)

		assert.ElementsMatch(t, []int{2}, listCurrentObjectIDs(h))
	})

	t.Run("only consider objects which are 'online' (id%2==0 from mock getter)", func(t *testing.T) {
		mockPers := persist.NewMockPersistence()
		mockGetter := get.MockObjectGetter{}
		h := NewObjectHandler(mockPers, &mockGetter, time.Second)

		h.Handle(2)
		h.Handle(3)
		assert.ElementsMatch(t, []int{2}, listCurrentObjectIDs(h))
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

func listCurrentObjectIDs(h *ObjectHandler) []int {
	var ids []int

	objs, _ := h.persistence.GetObjects()

	for objectID := range objs {
		ids = append(ids, objectID)
	}

	return ids
}
