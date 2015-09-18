package query

import (
  "strings"
  "io/ioutil"
  "golang.org/x/text/transform"
  "golang.org/x/text/encoding/charmap"
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
