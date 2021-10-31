package go_memory_test

import (
	"testing"

	"github.com/rihter007/go-memory"
)

func TestShared(t *testing.T) {
	var internal int = 10
	checkInternalObject := func(obj interface{}) {
		unpacked, ok := obj.(int)
		if !ok {
			t.Error("Failed to cast internal object to int")
		}
		if unpacked != 10 {
			t.Errorf("Result object is %d, expected %d", unpacked, internal)
		}
	}

	var releasedCount int
	shared := go_memory.NewShared(internal, func(obj interface{}) {
		checkInternalObject(obj)
		releasedCount++
	})
	checkInternalObject(shared.Get())
	if shared.Released() {
		t.Error("Released() of a newly created object returned true")
	}

	shared2 := shared
	checkInternalObject(shared2.Get())

	shared3 := shared2.AddRef()
	checkInternalObject(shared3.Get())

	if err := shared2.Close(); err != nil {
		t.Errorf("Error %v on object close", err)
	}
	if !shared2.Released() {
		t.Error("Released() of closed object should return true")
	}
	if !shared.Released() {
		t.Error("Released() of closed object should return true")
	}
	if err := shared.Close(); err != nil {
		t.Errorf("Error %v on object close", err)
	}

	if releasedCount != 0 {
		t.Errorf("object was released despite having references")
	}

	checkInternalObject(shared3.Get())
	if shared3.Released() {
		t.Errorf("Released() of not closed object shuld return false")
	}
	if err := shared3.Close(); err != nil {
		t.Errorf("Error %v on object close", err)
	}
	if !shared3.Released() {
		t.Errorf("Released() of closed object shuld return true")
	}
	if releasedCount != 1 {
		t.Errorf("released called was called %d number of times, expected %d", releasedCount, 1)
	}
}
