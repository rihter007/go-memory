package storage

import (
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

// Storage represents a memory-mapped storage
type Storage struct {
	directory string
	prefix    string
}

func New(directory string, prefix string) (*Storage, error) {
	// check access rights by creating a temporary file
	tmpFile, err := os.CreateTemp(directory, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create file in directory '%s', err: %w", directory, err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	if err := os.Remove(tmpFile.Name()); err != nil {
		return nil, err
	}

	return &Storage{
		directory: directory,
		prefix:    prefix,
	}, nil
}

func (s *Storage) Add(data []byte) (MemoryRegion, error) {
	if len(data) == 0 {
		return emptyMemoryRegion(), nil
	}
	f, err := os.CreateTemp(s.directory, s.prefix)
	if err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		// After the mmap() call has returned, the file descriptor, fd, can
		// be closed immediately without invalidating the mapping.
		f.Close()
		os.Remove(f.Name())
	}()
	if _, err := f.Write(data); err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to write data to file: %w", err)
	}
	mapping, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to map file into memory")
	}
	return newMemoryRegion(mapping), nil
}
