package mapping

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestTempMappings(t *testing.T) {
	inputBlobs := [][]byte{
		{0x1, 0x2, 0x3, 0x4},
		{0x42},
		{},
		nil,
	}
	var resultRegions []MemoryRegion
	for _, input := range inputBlobs {
		region, err := MapTemp(input, "")
		if err != nil {
			t.Fatalf("Failed to add memory object: '%v'", err)
		}
		if region.IsReleased() {
			t.Fatalf("Newly created region should not be released")
		}
		resultRegions = append(resultRegions, region)
	}

	for i := 0; i < len(inputBlobs); i++ {
		data := resultRegions[i].GetBytesUnsafe()
		if len(inputBlobs[i]) == 0 {
			if len(data) != 0 {
				t.Fatalf("Expected empty data, got [%X] bytes", data)
			}
		} else if !bytes.Equal(inputBlobs[i], data) {
			t.Fatalf("Expected [%X], got [%X] bytes", inputBlobs[i], data)
		}

		for k := 0; k != 2; k++ {
			newRegion := resultRegions[i]
			if err := newRegion.Close(); err != nil {
				t.Errorf("Failed to close memory region: %v", err)
			}
			if !newRegion.IsReleased() {
				t.Fatalf("Closed region should be released")
			}
			if err := resultRegions[i].Close(); err != nil {
				t.Errorf("Failed to close memory region: %v", err)
			}
			if !resultRegions[i].IsReleased() {
				t.Fatalf("Closed region should be released")
			}
		}
	}
}

func TestWritableTempMapping(t *testing.T) {
	mr, err := MapTemp([]byte{0x1, 0x2, 0x3, 0x4}, "")
	if err != nil {
		t.Fatalf("Failed to add memory object: '%v'", err)
	}
	mr2 := mr.AddRef()
	mr.GetBytesUnsafe()[0] = 0x42
	if err := mr.Close(); err != nil {
		t.Errorf("Failed to close memory region: %v", err)
	}

	expectedBytes := []byte{0x42, 0x2, 0x3, 0x4}
	if !bytes.Equal(expectedBytes, mr2.GetBytesUnsafe()) {
		t.Errorf("Expected [%X], got [%X]", expectedBytes, mr2.GetBytesUnsafe())
	}
}

func TestMapPath(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create temp file")
	}
	if _, err := f.Write([]byte{0x1, 0x2, 0x3, 0x4}); err != nil {
		t.Fatalf("Failed to write into file, err: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("Failed to close temp file: %v", err)
		}
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("Failed to remove temp file '%s', err: %v", f.Name(), err)
		}
	}()

	mr, err := MapPath(f.Name(), false)
	if err != nil {
		t.Fatalf("Failed to map path '%s', err: %v", f.Name(), err)
	}
	mr.GetBytesUnsafe()[0] = 0x42

	expected := []byte{0x42, 0x2, 0x3, 0x4}
	if !bytes.Equal(expected, mr.GetBytesUnsafe()) {
		t.Errorf("Expected [%X], actual [%X]", expected, mr.GetBytesUnsafe())
	}

	if err := mr.Close(); err != nil {
		t.Fatalf("Failed to close memory region, err: %v", err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("Failed to seek file, err: %v", err)
	}

	fileBody, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read file, err: %v", err)
	}

	if !bytes.Equal(expected, fileBody) {
		t.Errorf("Expected [%X], actual [%X]", expected, mr.GetBytesUnsafe())
	}
}
