package main

import (
	"cron/conf"
	"cron/routers"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
)

func main() {
	// 必须要先声明defer，否则不能捕获到panic异常
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("统一接受错误：", err)
		}
	}()

	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	config := conf.ReadConfig()
	go func() {
		http.ListenAndServe(config.Port, n)
	}()
	select {}

}
