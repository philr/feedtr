package feedtr

import (
  "errors"
  "fmt"
  "github.com/hashicorp/errwrap"
  "github.com/pkg/math"
  "log"
  "net/http"
  "sync"
  "time"
)

const fetchLimit = 100

// Saves the response to the cache.
func saveResponse(response *http.Response, cacheEntry CacheEntry) error {
  lastModifiedHeader := response.Header.Get("Last-Modified")

  if lastModifiedHeader == "" {
    lastModifiedHeader = response.Header.Get("Date")

    if lastModifiedHeader == "" {
      return errors.New("Missing Last-Modified and Date headers")
    }
  }

  lastModified, err := http.ParseTime(lastModifiedHeader)

  if err != nil {
    return err
  }

  return cacheEntry.Write(response.Body, lastModified)
}

// Fetches a source URL and saves the response to the cache.
func fetch(source string, cache Cache) (err error) {
  log.Printf("Fetching %s", source)

  request, err := http.NewRequest("GET", source, nil)
  if err != nil {
    return
  }

  request.Header.Add("User-Agent", "FeedTransformer/1 (https://github.com/philr/feedtransformer)")

  cacheEntry := cache.Entry(source)
  lastModified := cacheEntry.LastModified()

  if lastModified != nil {
    request.Header.Add("If-Modified-Since", lastModified.UTC().Format(time.RFC1123))
  }

  client := &http.Client{}
  response, err := client.Do(request)
  if err != nil {
    return
  }

  defer response.Body.Close()

  log.Printf("Fetched %s, got status %s", source, response.Status)

  if response.StatusCode == 200 {
    return saveResponse(response, cacheEntry)
  } else if response.StatusCode != 304 {
    return fmt.Errorf("Unexpected status %s", response.Status)
  }

  return
}

// Fetches all the sources specified by config, storing them in the given cache.
func FetchSources(config *Config, cache Cache) []error {
  sources := config.Sources()
  sc := make(chan string)
  ec := make(chan error)

  var errors []error

  go func() {
    for err := range ec {
      errors = append(errors, err)
    }
  }()

  var wg sync.WaitGroup
  wg.Add(len(sources))

  // Concurrently fetch, up to a maximum concurrency limit.
  for i := 0; i < math.MinInt(fetchLimit, len(sources)); i++ {
    go func() {
      for source := range sc {
        defer wg.Done()
        err := fetch(source, cache)
        if err != nil {
          ec <- errwrap.Wrapf(fmt.Sprintf("Error fetching %s: {{err}}", source), err)
        }
      }
    }()
  }

  go func() {
    for _, source := range sources {
      sc <- source
    }

    close(sc)
  }()

  wg.Wait()
  close(ec)

  return errors
}
