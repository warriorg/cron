package main

import (
	"cron/conf"
	"cron/models"
	"cron/routers"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	log.Println("num cpu", runtime.NumCPU())

	// 必须要先声明defer，否则不能捕获到panic异常
	defer func() {
		if err := recover(); err != nil {
			log.Println("统一接受错误：", err)
		}
	}()

	defer models.CloseDB()
 
	routers.SetupRoutes()
	// n := negroni.Classic()
	// n.UseHandler(router

	config := conf.ReadConfig()
	go func() {
		log.Println(http.ListenAndServe(config.Port, nil))
	}()
	select {}

}
