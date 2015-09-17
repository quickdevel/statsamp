package query

import (
  "strings"
  "io/ioutil"
  "code.google.com/p/go.text/transform"
  "code.google.com/p/go.text/encoding/charmap"
)

func transformToUTF8(source string) (string, error) {
  sourceReader := strings.NewReader(source)
  transformReader := transform.NewReader(sourceReader, charmap.Windows1251.NewDecoder())
  buf, err := ioutil.ReadAll(transformReader)
  if err != err {
    return "", err
  }
  return string(buf), nil
}
