package lib

import "os"

// FileCreateRemover ...
type FileCreateRemover interface {
	Create(name string) (*os.File, error)
	Remove(path string) error
}
