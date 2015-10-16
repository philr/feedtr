package feedtr

import (
  "fmt"
  "github.com/hashicorp/errwrap"
  "github.com/pkg/math"
  "log"
  "sync"
)

const processLimit = 100

// Runs the transformations for a single Output.
func processOutput(output Output, cache Cache, transforms Transforms, outputs Outputs) error {
  log.Printf("Processing %s", output.Name)

  outputEntry := outputs.Entry(output)
  outputLastModified := outputEntry.LastModified()
  cacheEntry := cache.Entry(output.Source)
  cacheEntryLastModified := cacheEntry.LastModified()

  if outputLastModified != nil && cacheEntryLastModified != nil &&
      (outputLastModified.After(*cacheEntryLastModified) || outputLastModified.Equal(*cacheEntryLastModified)) {
    log.Printf("Source hasn't changed for %s since the last transformation", output.Name)
    return nil
  }

  content, err := cacheEntry.Read()

  if err != nil {
    return err
  }

  for _, transformName := range output.Transforms {
    transform, err := transforms.Get(transformName)

    log.Printf("Running transform %s for %s", transformName, output.Name)

    if err != nil {
      return err
    }

    defer transform.Close()

    content, err = transform.Process(content)

    if err != nil {
      return err
    }
  }

  log.Printf("Writing output for %s", output.Name)

  outputEntry.Save(content)

  return nil
}

// Runs all the transforms for each Output specified in config. Sources are read
// from the given cache. Transforms are loaded from transforms. Transformed
// outputs are written to outputs.
func Process(config *Config, cache Cache, transforms Transforms, outputs Outputs) []error {
  oc := make(chan Output)
  ec := make(chan error)

  var errors []error

  go func() {
    for err := range ec {
      errors = append(errors, err)
    }
  }()

  var wg sync.WaitGroup
  wg.Add(len(config.Outputs))

  // Process outputs, up to a maximum concurrency limit.
  for i := 0; i < math.MinInt(processLimit, len(config.Outputs)); i++ {
    go func() {
      for output := range oc {
        defer wg.Done()
        err := processOutput(output, cache, transforms, outputs)
        if err != nil {
          ec <- errwrap.Wrapf(fmt.Sprintf("Error processing %s: {{err}}", output.Name), err)
        }
      }
    }()
  }

  go func() {
    for _, output := range config.Outputs {
      oc <- output
    }

    close(oc)
  }()

  wg.Wait()
  close(ec)

  return errors
}
