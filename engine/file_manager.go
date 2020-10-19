package engine

import "os"

// FileManger ...
type FileManger struct{}

// Create ...
func (fm FileManger) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Remove ...
func (fm FileManger) Remove(path string) error {
	return os.RemoveAll(path)
}
