package feedtr

import (
  "path/filepath"
)

// A source that can be used to obtain Transform instances.
type Transforms interface {
  // Gets a Transform by name.
  Get(name string) (Transform, error)
}

// An implementation of the Transforms instance using the file system to
// store transformation files.
type FileTransforms struct {
  path string
}

// Creates an instance of FileTransforms, to load transformation files from the
// given path.
func NewFileTransforms(path string) *FileTransforms {
  c := new(FileTransforms)
  c.path = path
  return c
}

// Gets a Transform by name.
func (transforms *FileTransforms) Get(name string) (Transform, error) {
  path := filepath.Join(transforms.path, name)
  return newXsltTransform(path)
}

