package config

import (
  "code.google.com/p/gcfg"
)

type Config struct {
  Main struct {
    Listen    string
  }
  DB struct {
    Host      string
    User      string
    Password  string
    Database  string
  }
  Updater struct {
    Version   string
  }
}

func GetConfigData() (Config, error) {
  var cfg Config
  err := gcfg.ReadFileInto(&cfg, "statsamp.cfg")
  return cfg, err
}
