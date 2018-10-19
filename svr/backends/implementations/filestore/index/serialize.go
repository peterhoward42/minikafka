package index

import (
	"encoding/gob"
	"fmt"
	"io"
)

// Encode is a serializer. It encodes the index into a byte stream and writes
// them to the output writer provided. See also the Decode sister method.
func (index *Index) Encode(writer io.Writer) error {
	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(index)
	if err != nil {
		return fmt.Errorf("encoder.Encode(): %v", err)
	}
	return nil
}

// Decode is a de-serializer. It populates the index by decoding the bytes
// read from the input reader provided. See also the Encode sister method.
func (index *Index) Decode(reader io.Reader) error {
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(index)
	if err != nil {
		return fmt.Errorf("decoder.Decode: %v", err)
	}
	return nil
}
