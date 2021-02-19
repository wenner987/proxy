package main

import (
	"log"
	"os"
	"path/filepath"
	"proxy/configure"
	"proxy/listen"
)

var appPath = ""
var configPath = "../conf/proxy.yaml"

func init() {
	log.SetFlags(log.Ldate|log.Lshortfile)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Get run path error! %v", err)
	}
	appPath = dir
	configPath = dir + "/" + configPath
	configure.ParseConfigure(configPath)
}

func main() {
	server := listen.Server{}
	server.Init()
	server.Listen()
}
