package mapping

import (
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

// MapFile maps specified file into memory
func MapFile(f *os.File, readOnly bool) (MemoryRegion, error) {
	var memProtection int
	if readOnly {
		memProtection = mmap.RDONLY
	} else {
		memProtection = mmap.RDWR
	}
	mapping, err := mmap.Map(f, memProtection, 0)
	if err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to map file into memory")
	}
	return newMemoryRegion(mapping), nil
}

// MapTemp puts data in a temporary file and returns result MemoryRegion
// If directory is the empty string, MapToTemp uses the default directory for temporary files, as returned by TempDir.
func MapTemp(data []byte, directory string) (MemoryRegion, error) {
	if len(data) == 0 {
		return emptyMemoryRegion(), nil
	}

	f, err := os.CreateTemp(directory, "memmapped")
	if err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to create temp file")
	}

	if _, err := f.Write(data); err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to write data into '%s', err: %w", f.Name(), err)
	}

	defer func() {
		// After the mmap() call has returned, the file descriptor, fd, can
		// be closed immediately without invalidating the mapping.
		f.Close()
		os.Remove(f.Name())
	}()
	return MapFile(f, false)
}

// MapPath maps file under specified path into memory
func MapPath(path string, readOnly bool) (MemoryRegion, error) {
	var flag int
	if readOnly {
		flag = os.O_RDONLY
	} else {
		flag = os.O_RDWR
	}
	f, err := os.OpenFile(path, flag, 0)
	if err != nil {
		return MemoryRegion{}, fmt.Errorf("failed to open file: '%s'", path)
	}
	return MapFile(f, readOnly)
}
