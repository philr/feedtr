package feedtr

import (
  "crypto/sha1"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "path/filepath"
  "time"
)

// A cache used to store fetched feed sources.
type Cache interface {
  // Returns a CacheEntry for a given feed source URL. A result will always be
  // returned regardless of whether the entry exists.
  Entry(url string) CacheEntry
}

// A entry for a feed source in a Cache.
type CacheEntry interface {
  // Writes content from the given Reader to the cache entry. The given
  // lastModified time is recorded along with the content.
  Write(reader io.Reader, lastModified time.Time) error

  // Returns the last modified Time of the entry or nil if nothing has been
  // stored for the entry.
  LastModified() *time.Time

  // Returns the cached content for the cache entry as a byte slice.
  Read() ([]byte, error)
}

// A implementation of the Cache interface using the file system to store
// entries.
type FileCache struct {
  path string
}

// A implementation of the CacheEntry interface referencing a file on the file
// system.
type FileCacheEntry struct {
  path string
}

// Creates and returns a new FileCache to write cache entries to the given path.
// The directory specified by the path is created if it does not exist.
func NewFileCache(path string) (*FileCache, error) {
  err := os.MkdirAll(path, 0700)

  if err != nil {
    return nil, err
  }

  c := new(FileCache)
  c.path = path

  return c, err
}

// Returns a CacheEntry for a given feed source URL. A result will always be
// returned regardless of whether the entry exists.
func (cache *FileCache) Entry(url string) CacheEntry {
  h := sha1.New()
  h.Write([]byte(url))
  hash := h.Sum(nil)

  result := new(FileCacheEntry)
  result.path = filepath.Join(cache.path, fmt.Sprintf("%x", hash))

  return result
}

// Writes content from the given Reader to the cache entry. The given
// lastModified time is recorded along with the content.
func (entry *FileCacheEntry) Write(reader io.Reader, lastModified time.Time) error {
  tempFile := entry.path + ".tmp"
  file, err := os.Create(tempFile)

  if err != nil {
    return err
  }

  _, err = io.Copy(file, reader)

  file.Close()

  if err != nil {
    return err
  }

  err = os.Rename(tempFile, entry.path)

  if err != nil {
    os.Remove(tempFile)
    return err
  }

  err = os.Chtimes(entry.path, time.Now(), lastModified)
  return err
}

// Returns the last modified Time of the entry or nil if nothing has been
// stored for the entry.
func (entry *FileCacheEntry) LastModified() *time.Time {
  return fileLastModified(entry.path)
}

// Returns the cached content for the cache entry as a byte slice.
func (entry *FileCacheEntry) Read() ([]byte, error) {
  return ioutil.ReadFile(entry.path)
}
