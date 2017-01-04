package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string `yaml:"port"`
	DB   struct {
		URL string `yaml:'url'`
	}
}

func ReadConfig() *Config {
	var conf *Config
	data, err := ioutil.ReadFile("conf/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return conf
}
