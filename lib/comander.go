package lib

import "io"

// Commander ...
type Commander interface {
	StdinPipe() (io.WriteCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Output() ([]byte, error)
}
