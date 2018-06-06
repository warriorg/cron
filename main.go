package main

import (
	"cron/conf"
	"cron/routers"
	"flag"
	"fmt"
	"ihome-web/models"
	"log"
	"net/http"
	"os"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	// log.Println("num cpu", runtime.NumCPU())
	cmd := parseCmd()
	// fmt.Println(cmd)
	if cmd.versionFlag {
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag {
		printUsage()
	} else {

		if cmd.modeFlag == "web" {
			startWebServer(cmd)
		} else if cmd.modeFlag == "grpc" {
			startGrpcServer(cmd)
		}

	}

}

func startGrpcServer(cmd *Cmd) {

}

func startWebServer(cmd *Cmd) {

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

// Cmd cmd
type Cmd struct {
	helpFlag    bool
	versionFlag bool
	configFlag  string
	modeFlag    string
	args        []string
}

func parseCmd() *Cmd {
	cmd := &Cmd{}
	flag.Usage = printUsage
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.StringVar(&cmd.modeFlag, "mode", "web", "mode web or grpc")
	flag.StringVar(&cmd.configFlag, "c", "./", "set configuration file")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		cmd.args = args
	}
	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-options] [--mode web|grpc][-c filename] [args...] \n", os.Args[0])
}
