package main

import (
	"assessment/pkg/get"
	"assessment/pkg/handle"
	"assessment/pkg/persist"
	"time"
)

func configureObjectHandler() *handle.ObjectHandler {
	var persistence persist.Persistor
	if !*useMockDB {
		persistence = persist.NewPostgres()
	} else {
		persistence = persist.NewMockPersistence()
	}

	err := persistence.Connect()
	if err != nil {
		panic(err)
	}

	var getter get.Getter
	if *objectSource != mockObjectSourceIdentifier {
		getter = get.NewRemoteObjectGetter(*objectSource)
	} else {
		getter = &get.MockObjectGetter{}
	}

	return handle.NewObjectHandler(persistence, getter, time.Second*time.Duration(*objectLifespan))
}
