package main

import (
	"cron/routers"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string `yaml:"port"`
}

func main() {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	config := new(Config)
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	http.ListenAndServe(config.Port, n)
}
