package mapping

import (
	"github.com/edsrzf/mmap-go"
	"github.com/rihter007/go-memory"
)

type MemoryRegion struct {
	shared memory.Shared
}

func (mr MemoryRegion) AddRef() MemoryRegion {
	return MemoryRegion{
		shared: mr.shared.AddRef(),
	}
}

func (mr MemoryRegion) GetBytesUnsafe() []byte {
	internal := mr.shared.Get()
	if internal == nil {
		return nil
	}
	return internal.(mmap.MMap)
}

func (mr MemoryRegion) IsReleased() bool {
	return mr.shared.Released()
}

func (mr MemoryRegion) Close() error {
	return mr.shared.Close()
}

func newMemoryRegion(mapping mmap.MMap) MemoryRegion {
	if mapping == nil {
		panic("mapping should not be nil")
	}
	return MemoryRegion{
		shared: memory.NewShared(mapping, func(obj interface{}) error {
			return mapping.Unmap()
		}),
	}
}

func emptyMemoryRegion() MemoryRegion {
	return MemoryRegion{
		shared: memory.NewShared(nil, nil),
	}
}
