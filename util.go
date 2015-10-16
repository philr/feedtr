package feedtr

import (
  "os"
  "time"
)

// Returns the last modified time of a file or nil if the file does not exist
// (or an error occurs accessing the file info).
func fileLastModified(path string) *time.Time {
  fileInfo, err := os.Stat(path)

  if err == nil {
    result := fileInfo.ModTime()
    return &result
  }

  return nil
}
