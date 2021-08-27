package main

import (
	"context"
	"testing"
	"time"
)

const objecLifespan = 1 * time.Second

func TestService(t *testing.T) {
	tester := newTester(t)
	srv := serve()

	// very basic post and expect
	tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// reset
	time.Sleep(objecLifespan)

	// expect table to be empty after object lifespan has passed
	tester.postAndExpect([]int{}, []int{})

	// post object id of object which is !online (because id%2 != 0) and expect it not to be in table
	tester.postAndExpect([]int{2, 3}, []int{2})

	// post another id while expecting the existing table entry (2) to persist (lifespan hasn't passed yet)
	tester.postAndExpect([]int{4}, []int{2, 4})

	// reset
	time.Sleep(objecLifespan)

	// basic post and expect
	tester.postAndExpect([]int{2, 4}, []int{2, 4})

	// wait half a lifespan
	time.Sleep(objecLifespan / 2)

	// post only `2` to refresh lifespan while `4` remains untouched
	tester.postAndExpect([]int{2}, []int{2, 4})

	// wait half a lifespan
	time.Sleep(objecLifespan / 2)

	// post nothing, expect `2` to still be there since only lifespan/2 has passed since last post of `2`
	tester.postAndExpect([]int{}, []int{2})

	// reset
	time.Sleep(objecLifespan)

	// post same id multiple times, but expect only one to be in table
	tester.postAndExpect([]int{2, 2, 2, 2}, []int{2})

	srv.Shutdown(context.Background())
}
