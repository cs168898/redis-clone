package main

import "io"

type Writer struct {
	writer io.Writer
}

// constructor function
func NewWriter(w io.Writer) *Writer {
	// we initialize the writer field of the writer struct with
	// the args passed in the function, therefore creating
	// an instance of this writer.
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	// the write method returns the number of bytes return and the error
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
	
}