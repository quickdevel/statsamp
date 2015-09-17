package serverlist

import (
  "fmt"
  "bufio"
  "net/http"
  "io"
)

type Request struct {
  Version   string
  Response  *http.Response
  Reader    *bufio.Reader
}

func (this *Request) Exec() error {
  url := fmt.Sprintf("http://lists.sa-mp.com/%s/servers", this.Version)
  response, err := http.Get(url)
  if err != nil {
    return err
  }
  reader := bufio.NewReader(response.Body)
  this.Response = response;
  this.Reader = reader;
  return nil
}

func (this *Request) ReadNext() (bool, string, error) {
  line, _, err := this.Reader.ReadLine()
  if err != nil {
    if err != io.EOF {
      return false, "", nil
    }
    return false, "", err
  }
  return true, string(line), nil
}

func (this *Request) ReadAll() ([]string, error) {
  var output []string
  for {
    line, _, err := this.Reader.ReadLine()
    if err != nil {
      if err != io.EOF {
        return output, err
      }
      return output, nil
    }
    output = append(output, string(line))
  }
  return output, nil
}

func NewRequest(version string) Request {
  var req Request
  req.Version = version
  return req
}
