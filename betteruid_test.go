package betteruid

import (
	"fmt"
	"sync"
	"testing"
)

const (
	//for manual debugging
	showIDs = true
)

func TestBasic(t *testing.T) {
	id1 := New()
	if len(id1) != 20 {
		t.Fatalf("len(id1) != 20 (=%d)", len(id1))
	}
	id2 := New()
	if len(id2) != 20 {
		t.Fatalf("len(id2) != 20 (=%d)", len(id2))
	}

	if id1 == id2 {
		t.Fatalf("generated same ids (%s,%s)", id1, id2)
	}
	if showIDs {
		fmt.Printf("%s\n", id1)
		fmt.Printf("%s\n", id2)
	}
}

func TestBatch(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go batch(t, &wg)
	}
	wg.Wait()
}

func batch(t *testing.T, wg *sync.WaitGroup) {
	ids := make(map[string]bool, 1000000)
	prev := ""
	for i := 0; i < 1000000; i++ {
		id := New()
		if _, exists := ids[id]; exists {
			t.Fatalf("generate duplicate id %s", id)
		}
		ids[id] = true
		if prev != "" {
			if id <= prev {
				t.Fatalf("id(%s) must > prev(%s)", id, prev)
			}
		}
		prev = id
	}
	wg.Done()
}
