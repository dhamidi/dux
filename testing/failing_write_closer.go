package testing

import (
	"fmt"
	"io"
)

// FailingWriteCloser wraps an io.WriteCloser and returns an error
// when trying to write to the underlying WriteCloser.
type FailingWriteCloser struct {
	io.WriteCloser
}

// NewFailingWriteCloser returns a new FailingWriteCloser that wraps wc.
func NewFailingWriteCloser(wc io.WriteCloser) *FailingWriteCloser {
	return &FailingWriteCloser{WriteCloser: wc}
}

// Write always returns an error and 0.
func (wc *FailingWriteCloser) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Writing %d bytes failed", len(p))
}
