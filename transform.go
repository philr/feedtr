package feedtr

import (
  "github.com/jbowtie/ratago/xslt"
  "github.com/ThomsonReutersEikon/gokogiri/xml"
  "io/ioutil"
  "os"
)

// A transformation that can be run on a feed.
type Transform interface {
  // Runs the transformation on the given input, returning an output as a byte
  // slice.
  Process(input []byte) ([]byte, error)

  // Closes the Transform. Must be called after the instance is no longer
  // needed.
  Close()
}

// A Transform based on an XSLT stylesheet.
type xsltTransform struct {
  stylesheet *xslt.Stylesheet
}

// Creates an xlstTransform from file specified by path.
func newXsltTransform(path string) (*xsltTransform, error) {
  doc, err := xml.ReadFile(path, xml.StrictParseOption)

  if err != nil {
    return nil, err
  }

  stylesheet, err := xslt.ParseStylesheet(doc, path)

  if err != nil {
    doc.Free()
    return nil, err
  }

  result := new(xsltTransform)
  result.stylesheet = stylesheet
  return result, nil
}

// Runs the transformation on the given input, returning an output as a byte
// slice.
func (transform *xsltTransform) Process(input []byte) ([]byte, error) {
  // gokogiri can automatically determine the encoding when parsing a
  // file, but not when loading from memory. Save to a temporary file until
  // a better method can be found.

  file, err := ioutil.TempFile(os.TempDir(), "ftxlstin")

  if err != nil {
    return nil, err
  }

  defer os.Remove(file.Name())

  _, err = file.Write(input)

  if err != nil {
    return nil, err
  }

  doc, err := xml.ReadFile(file.Name(), xml.StrictParseOption)

  if err != nil {
    return nil, err
  }

  defer doc.Free()

  options := xslt.StylesheetOptions{false, nil}
  output, err := transform.stylesheet.Process(doc, options)

  if err != nil {
    return nil, err
  }

  return []byte(output), nil
}

// Closes the Transform. Must be called after the instance is no longer needed.
func (transform *xsltTransform) Close() {
  transform.stylesheet.Doc.Free()
}
