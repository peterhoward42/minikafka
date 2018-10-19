package index

import (
    "os"
    "fmt"
)

// SaveIndex serializes the given index into a byte stream representation,
// and saves this as a binary file.
func SaveIndex(index *Index, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("os.Create(): %v", err)
	}
	defer file.Close()
	err = index.Encode(file)
	if err != nil {
		return fmt.Errorf("Encode(): %v", err)
	}
	return nil
}

// RetrieveIndexFromDisk reads the bytes from the nominated file which was created
// using the SaveIndex sister method, and deserializes them to make an 
// Index object.
func RetrieveIndexFromDisk(filepath string) (*Index, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	index := NewIndex()
	err = index.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Decode(): %v", err)
	}
	return index, nil
}
