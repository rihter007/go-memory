package storage

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestStorage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestStorage")
	if err != nil {
		t.Fatalf("Failed to create temp directory: '%v'", err)
	}
	storage, err := New(tempDir, "index_")
	if err != nil {
		t.Fatalf("Failed to create storage: '%v'", err)
	}
	if storage == nil {
		t.Fatal("Result storage is nil")
	}

	inputBlobs := [][]byte{
		{0x1, 0x2, 0x3, 0x4},
		{0x42},
		{},
		nil,
	}
	var resultRegions []MemoryRegion
	for _, input := range inputBlobs {
		region, err := storage.Add(input)
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
