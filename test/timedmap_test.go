package timedmap

import (
	"fmt"
	"testing"
	"time"

	"github.com/codenamoo/timedmap"
)

func TestFunction(t *testing.T) {
	m := timedmap.NewTimedMap()

	// put new key:value
	fmt.Printf("timedmap demo\n")
	fmt.Printf("put new key:value with TTL 3\n\n")
	m.PutRaw("key", "value", 3)

	fmt.Printf("sleep 5secs\n")
	fmt.Printf("then key will be expired\n")
	time.Sleep(5 * time.Second)

	// key probably expired
	if v, err := m.Get("key"); err != nil {
		fmt.Printf("ERROR: %s\n\n", err.Error())
	} else {
		// huh?
		fmt.Printf("[key]: %v", v)
		t.FailNow()
	}

	// put new key:value again
	fmt.Printf("put new key:value again with TTL 10\n")
	m.PutRaw("key", "value", 10)

	if ttl, err := m.GetTTL("key"); err != nil {
		// huh?
		fmt.Printf("ERROR: %s\n", err.Error())
		t.FailNow()
	} else {
		// expect 10
		fmt.Printf("I expect TTL 10\n")
		fmt.Printf("KEY: %s, TTL: %d\n\n", "key", ttl)
	}

	fmt.Printf("sleep 5secs\n")
	time.Sleep(5 * time.Second)

	if ttl, err := m.GetTTL("key"); err != nil {
		// huh?
		fmt.Printf("ERROR: %s\n", err.Error())
		t.FailNow()
	} else {
		// expect 5
		if ttl != 5 {
			// huh? timing issue?
		}
		fmt.Printf("I expect TTL 5\n")
		fmt.Printf("KEY: %s, TTL: %d\n", "key", ttl)
	}

	// reset ttl
	fmt.Printf("Reset TTL\n\n")
	m.Touch("key")

	fmt.Printf("sleep 6secs\n")
	time.Sleep(6 * time.Second)

	// expect 4
	if ttl, err := m.GetTTL("key"); err != nil {
		// huh?
		fmt.Printf("ERROR: %s\n", err.Error())
		t.FailNow()
	} else {
		fmt.Printf("I expect TTL 4\n")
		fmt.Printf("KEY: %s, TTL: %d\n", "key", ttl)
	}
}
