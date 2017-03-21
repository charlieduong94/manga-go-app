package config

import (
  "sync"
  "io/ioutil"
  "github.com/golang/glog"

  "gopkg.in/yaml.v2"
)

type Config struct {
  Env string `yaml:"env"`
  DynamoDBEndpoint string `yaml:"dynamoDBEndpoint"`
}

var config Config
var once sync.Once

func GetConfig () Config {
  once.Do(func () {
    glog.Info("Loading config")

    rawConfig, err := ioutil.ReadFile("config.yml")

    // TODO: handle errors properly and don't panic
    if err != nil {
      panic(err)
    }

    err = yaml.Unmarshal(rawConfig, &config)
    if err != nil {
      panic(err)
    }

    glog.Infof("Config: %+v", config)
  })

  return config
}
