package json

import (
	"io"
)

// An Encoder writes JSON values to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the JSON encoding of v to the stream, followed by a newline character.
//
// See the documentation for Marshal for details about the conversion of Go values to JSON.
func (enc *Encoder) Encode(v interface{}) error {
	// Marshal the value
	data, err := Marshal(v)
	if err != nil {
		return err
	}

	// Write to the stream
	if _, err := enc.w.Write(data); err != nil {
		return err
	}

	// Write newline
	if _, err := enc.w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}
