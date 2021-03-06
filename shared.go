package memory

import (
	"sync/atomic"
)

// Shared represents a shared object
// We explicitly disable an ability for external clients to unref the object
type Shared struct {
	*shared
}

type shared struct {
	obj       interface{}
	onRelease func(obj interface{}) error
	refCount  *int32
	closed    uint32
}

func NewShared(obj interface{}, onRelease func(obj interface{}) error) Shared {
	refCount := int32(1)
	s := Shared{new(shared)}
	s.obj = obj
	s.onRelease = onRelease
	s.refCount = &refCount
	return s
}

func (s Shared) Get() interface{} {
	if s.Released() {
		panic("Attempt to get object after it was released")
	}
	return s.obj
}

func (s Shared) AddRef() Shared {
	if s.Released() {
		panic("Attempt to add reference to object after it was released")
	}
	atomic.AddInt32(s.refCount, 1)
	result := Shared{new(shared)}
	result.obj = s.obj
	result.onRelease = s.onRelease
	result.refCount = s.refCount
	return result
}

func (s Shared) Close() error {
	if !atomic.CompareAndSwapUint32(&s.closed, 0, 1) {
		return nil
	}
	newCount := atomic.AddInt32(s.refCount, -1)
	if newCount == 0 && s.onRelease != nil {
		return s.onRelease(s.obj)
	}
	return nil
}

func (s *shared) Released() bool {
	return atomic.LoadUint32(&s.closed) == 1
}
