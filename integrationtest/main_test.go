package main

import (
	"assessment/pkg/server"
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

const OBJECT_LIFESPAN = 1
const OBJECT_SERVER_PORT = ":9010"
const SERVICE_PORT = ":9090"

func TestService(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error reading env file: %s", err)
	}

	killPostgres, err := startPostgresContainer()
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
		ObjectLifespan: OBJECT_LIFESPAN,
		ObjectSource:   "http://localhost" + OBJECT_SERVER_PORT,
	}

	ctx := context.Background()

	srv := server.Serve(SERVICE_PORT, c)
	defer srv.Shutdown(ctx)
	objectServer := serveObjectSource()
	defer objectServer.Shutdown(ctx)
	time.Sleep(OBJECT_LIFESPAN * time.Second)

	// very basic post and expect
	tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// reset
	time.Sleep(OBJECT_LIFESPAN * time.Second)

	// // expect table to be empty after object lifespan has passed
	// tester.postAndExpect([]int{}, []int{})

	// // post object id of object which is !online (because id%2 != 0) and expect it not to be in table
	// tester.postAndExpect([]int{2, 3}, []int{2})

	// // post another id while expecting the existing table entry (2) to persist (lifespan hasn't passed yet)
	// tester.postAndExpect([]int{4}, []int{2, 4})

	// // reset
	// time.Sleep(OBJECT_LIFESPAN * time.Second)

	// // basic post and expect
	// tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// // wait half a lifespan
	// time.Sleep(OBJECT_LIFESPAN * time.Second / 2)

	// // post only `2` to refresh lifespan while `4` remains untouched
	// tester.postAndExpect([]int{2}, []int{2, 4})

	// // wait half a lifespan
	// time.Sleep(OBJECT_LIFESPAN * time.Second / 2)

	// // post nothing, expect `2` to still be there since only lifespan/2 has passed since last post of `2`
	// tester.postAndExpect([]int{}, []int{2})

	// // reset
	// time.Sleep(OBJECT_LIFESPAN * time.Second)

	// // post same id multiple times, but expect only one to be in table
	// tester.postAndExpect([]int{2, 2, 2, 2}, []int{2})
}
