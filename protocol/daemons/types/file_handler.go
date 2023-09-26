package types

import (
	"os"
)

// FileHandlerImpl is the struct that implements the `FileHandler` interface.
type FileHandlerImpl struct{}

// Ensure the `FileHandlerImpl` struct is implemented at compile time.
var _ FileHandler = (*FileHandlerImpl)(nil)

// FileHandler is an interface that encapsulates the os function `RemoveAll`.
type FileHandler interface {
	RemoveAll(path string) error
}

// RemoveAll wraps `os.RemoveAll` which removes everything at a given path.
func (r *FileHandlerImpl) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
