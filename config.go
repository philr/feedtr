package feedtr

// Specifies the outputs to produce.
type Config struct {
  // A list of the outputs to be produced, each specified as an Output.
  Outputs []Output
}

// An configured output to produce.
type Output struct {
  // Name of the output to be written (typically a file name).
  Name       string

  // URL of the source feed to be fetched an transformed.
  Source     string

  // List of transforms to be applied to the feed (i.e. file names of the
  // transform files).
  Transforms []string
}

// Returns all the source URLs referenced in Outputs.
func (config *Config) Sources() []string {
  sources := make(map[string]bool)

  for _, output := range config.Outputs {
    sources[output.Source] = true
  }

  result := make([]string, len(sources))

  i := 0
  for source, _ := range sources {
    result[i] = source
    i++
  }

  return result
}
