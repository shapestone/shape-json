package json

import (
	"io"
)

// A Decoder reads and decodes JSON values from an input stream.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read data from r
// beyond the JSON values requested.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next JSON-encoded value from its input and stores it
// in the value pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion
// of JSON into a Go value.
func (dec *Decoder) Decode(v interface{}) error {
	// Use ParseReader to parse JSON from the stream
	node, err := ParseReader(dec.r)
	if err != nil {
		return err
	}

	// Use the same unmarshal logic as Unmarshal
	return unmarshalFromNode(node, v)
}
