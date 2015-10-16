package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/philr/feedtr"
  "io"
  "os"
  "path/filepath"
)

// Writes an error and exits if e is not nil.
func check(e error) {
  if e != nil {
    fmt.Fprintln(os.Stderr, e.Error())
    os.Exit(1)
  }
}

func checkAll(e []error) {
  if len(e) > 0 {
    for _, e := range e {
      fmt.Fprintln(os.Stderr, e.Error())
    }

    os.Exit(1)
  }
}

// Loads a JSON representation of Config from the given Reader.
func loadConfigFrom(reader io.Reader) (result *feedtr.Config, err error) {
  decoder := json.NewDecoder(reader)
  result = new(feedtr.Config)
  err = decoder.Decode(result)
  return
}

// Loads a JSON representation of Config from either standard input (if
// configFile is empty) or a file (if configFile is set). When loading from a
// file the current directory is changed to the directory containing
// configFile.
func loadConfig(configFile string) (result *feedtr.Config, err error) {
  if configFile == "" {
    return loadConfigFrom(os.Stdin)
  } else {
    var file *os.File
    file, err = os.Open(configFile)

    if err != nil {
      return
    }

    defer file.Close()
    os.Chdir(filepath.Dir(configFile))
    return loadConfigFrom(file)
  }
}

// Loads the Config, fetches sources, then runs transformations to create
// outputs.
func main() {
  var configFile string
  flag.StringVar(&configFile, "config-file", "", "Path to config file from")
  flag.Parse()

  config, err := loadConfig(configFile)
  check(err)

  cache, err := feedtr.NewFileCache("cache")
  check(err)

  errs := feedtr.FetchSources(config, cache)
  checkAll(errs)

  transforms := feedtr.NewFileTransforms("transforms")
  outputs, err := feedtr.NewFileOutputs("outputs")
  check(err)

  errs = feedtr.Process(config, cache, transforms, outputs)
  checkAll(errs)
}
