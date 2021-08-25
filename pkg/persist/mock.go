package persist

import (
	"assessment/pkg/object"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type MockPersistence struct {
	objects map[int]object.Object
	mu      sync.Mutex
}

func NewMockPersistence() *MockPersistence {
	m := MockPersistence{
		objects: make(map[int]object.Object),
		mu:      sync.Mutex{},
	}

	go m.keepPrinting()

	return &m
}

func (m *MockPersistence) keepPrinting() {
	ticker := time.NewTicker(time.Second)

	for {
		<-ticker.C
		b, err := json.MarshalIndent(m.objects, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(string(b))
		fmt.Println("------------------------------")
	}
}

func (m *MockPersistence) WriteObject(o object.Object) error {
	m.mu.Lock()
	m.objects[o.ObjectID] = o
	m.mu.Unlock()
	return nil
}

func (m *MockPersistence) DeleteObject(objectID int, lastSeen int64) error {
	m.mu.Lock()
	for _, o := range m.objects {
		if o.ObjectID == objectID && o.LastSeen == lastSeen {
			delete(m.objects, o.ObjectID)
		}
	}
	m.mu.Unlock()
	return nil
}
