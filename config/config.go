package config

import (
  "sync"
  "io/ioutil"
  "github.com/golang/glog"

  "gopkg.in/yaml.v2"
)

type Config struct {
  Env string `yaml:"env"`
  DynamoDbEndpoint string `yaml:"dynamoDbEndpoint"`
}

var configLoaded = false
var config Config
var once sync.Once

func Load (configFile string) {
  once.Do(func () {
    glog.Info("Loading config")

    rawConfig, err := ioutil.ReadFile(configFile)

    // TODO: handle errors properly and don't panic
    if err != nil {
      glog.Info("Could not read config file..")
    } else {
      err = yaml.Unmarshal(rawConfig, &config)
      if err != nil {
        panic(err)
      }
    }

    glog.Infof("Config: %+v", config)
    configLoaded = true
  })
}

func GetConfig () Config {
  if configLoaded == false {
    panic("Config not loaded yet")
  }
  return config
}
