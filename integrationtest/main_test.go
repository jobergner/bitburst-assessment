package integrationtest

import (
	"assessment/integrationtest/container"
	"assessment/pkg/object"
	"assessment/pkg/server"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	objectLifespan = 1 * time.Second
)

func runTestCases(tester *integrationTester) {
	// very basic post and expect
	tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// reset
	time.Sleep(objectLifespan)

	// expect table to be empty after object lifespan has passed
	tester.postAndExpect([]int{}, []int{})

	// post object id of object which is !online (because id%2 != 0) and expect it not to be in table
	tester.postAndExpect([]int{2, 3}, []int{2})

	// post another id while expecting the existing table entry (2) to persist (lifespan hasn't passed yet)
	tester.postAndExpect([]int{4}, []int{2, 4})

	// reset
	time.Sleep(objectLifespan)

	// basic post and expect
	tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// wait half a lifespan
	time.Sleep(objectLifespan / 2)

	// post only `2` to refresh lifespan while `4` remains untouched
	tester.postAndExpect([]int{2}, []int{2, 4})

	// wait half a lifespan
	time.Sleep(objectLifespan / 2)

	// post nothing, expect `2` to still be there since only lifespan/2 has passed since last post of `2`
	tester.postAndExpect([]int{}, []int{2})

	// reset
	time.Sleep(objectLifespan)

	// post same id multiple times, but expect only one to be in table
	tester.postAndExpect([]int{2, 2, 2, 2}, []int{2})
}

func TestService(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error reading env file: %s", err)
	}

	killPostgres, err := container.StartPostgres()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer killPostgres()

	tester, err := newTester(t)
	if err != nil {
		t.Fatalf(err.Error())
	}

	c := server.ServerConfig{
		UseMockDB:      false,
		ObjectLifespan: objectLifespan,
		ObjectSource:   "http://localhost" + objectServerPort,
	}

	ctx := context.Background()

	fillDB(tester)
	srv := server.Serve(servicePort, c)
	expectExpiredObjectsDeleted(tester)

	isRunning := waitForIt(fmt.Sprintf("http://localhost%s", servicePort))
	if !isRunning {
		t.Fatalf("could not connect to service")
	}
	defer srv.Shutdown(ctx)

	objectServer := serveObjectSource()
	isRunning = waitForIt(fmt.Sprintf("http://localhost%s", objectServerPort))
	if !isRunning {
		t.Fatalf("could not connect to object server")
	}
	defer objectServer.Shutdown(ctx)

	runTestCases(tester)
}

// check if the expired object got deleted after server startup
func expectExpiredObjectsDeleted(tester *integrationTester) {
	objs, err := tester.pg.GetObjects()
	if err != nil {
		tester.t.Fatalf(err.Error())
	}

	_, ok := objs[1]
	assert.True(tester.t, ok)

	_, ok = objs[2]
	assert.False(tester.t, ok)

	// wait and clear db before running next test cases
	time.Sleep(objectLifespan)
	tester.pg.DeleteObjectsOlderThan(objectLifespan)
}

// fillDB fills the postgres with some objects to check if initializing the server later removes expired objects on startup
func fillDB(tester *integrationTester) {
	object1 := object.Object{
		ObjectID: 1,
		Online:   true,
		LastSeen: time.Now().UnixNano(),
	}
	tester.pg.WriteObject(object1)

	object2 := object.Object{
		ObjectID: 2,
		Online:   true,
		// setting an expired time
		LastSeen: time.Now().Add(objectLifespan * -1).UnixNano(),
	}
	tester.pg.WriteObject(object2)

	// wait for write
	time.Sleep(time.Millisecond * 100)
}
