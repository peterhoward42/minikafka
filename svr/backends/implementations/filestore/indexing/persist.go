package indexing

import (
    "os"
    "fmt"
)

// Save serializes the index into a byte stream representation, and saves this 
// as a binary file.
func (index *Index) Save(filepath string) error {
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

// PopulateFromDisk reads the bytes from the nominated file which was created
// using the SaveIndex sister method, and deserializes them popualate this
// Index object.
func (index *Index) PopulateFromDisk(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("os.Open(): %v", err)
	}
	defer file.Close()
	err = index.Decode(file)
	if err != nil {
		return fmt.Errorf("Decode(): %v", err)
	}
	return nil
}
