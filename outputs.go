package feedtr

import (
  "io/ioutil"
  "os"
  "path/filepath"
  "time"
)

// The output location for transformed feeds.
type Outputs interface {
  // Returns an entry for an individual configured transformed feed output.
  Entry(output Output) OutputEntry
}

// An entry for a transformed feed output.
type OutputEntry interface {
  // Stores the given content.
  Save(content []byte) error

  // Returns the last modified time of the output or nil if the entry has not
  // yet been written.
  LastModified() *time.Time
}

// An implementation of the Outputs interface using the file system for storage.
type FileOutputs struct {
  path string
}

// An implementation of the OutputEntry interface using the file system for
// storage.
type FileOutputEntry struct {
  path string
}

// Returns a new FileOutputs instance to write outputs to the given path. The
// directory referenced by the path is created if it does not exist.
func NewFileOutputs(path string) (*FileOutputs, error) {
  err := os.MkdirAll(path, 0755)

  if err != nil {
    return nil, err
  }

  c := new(FileOutputs)
  c.path = path
  return c, nil
}

// Returns an entry for an individual configured transformed feed output.
func (outputs *FileOutputs) Entry(output Output) OutputEntry {
  result := new(FileOutputEntry)
  result.path = filepath.Join(outputs.path, output.Name)

  return result
}

// Stores the given content.
func (entry *FileOutputEntry) Save(content []byte) error {
  tempPath := entry.path + ".tmp"

  err := ioutil.WriteFile(tempPath, content, 0644)

  if err != nil {
    return err
  }

  // The umask is applied to the FileMode passed to WriteFile. Ensure the
  // permissions get set to allow others read access (since it is likely the
  // outputs will be published on a web server).
  err = os.Chmod(tempPath, 0644)

  if err != nil {
    os.Remove(tempPath)
    return err
  }

  err = os.Rename(tempPath, entry.path)

  if err != nil {
    os.Remove(tempPath)
    return err
  }

  return nil
}

// Returns the last modified time of the output or nil if the entry has not
// yet been written.
func (entry *FileOutputEntry) LastModified() *time.Time {
  return fileLastModified(entry.path)
}
