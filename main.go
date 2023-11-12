package main

import (
	"fmt"
	"github.com/dwnGnL/ddos-pow/cmd"
	"github.com/dwnGnL/ddos-pow/config"
	"log"
	"os"

	"github.com/dwnGnL/ddos-pow/lib/goerrors"

	"github.com/sirupsen/logrus"
)

const (
	CLIENT = "client"
	SERVER = "server"
)

var Version = "v0.0.1"

func main() {
	args := os.Args

	if len(args) != 2 {
		log.Fatalf("service must be runned with one parameter!")
		return
	}

	os.Setenv("config", "./config.yaml")

	switch args[1] {
	case CLIENT:
		cfg := config.FromFile(os.Getenv("config"))
		fmt.Println(cfg)
		intLogger(cfg.LogLevel)
		err := cmd.StartClient(cfg)
		if err != nil {
			log.Fatalf("app run: %s", err)
		}
	case SERVER:
		cfg := config.FromFile(os.Getenv("config"))
		fmt.Println(cfg)
		intLogger(cfg.LogLevel)
		err := cmd.StartServer(cfg)
		if err != nil {
			log.Fatalf("app run: %s", err)
		}
	}
}

func intLogger(logLevel string) {
	var formatter logrus.Formatter = new(logrus.JSONFormatter)
	if os.Getenv("LOG_FORMAT") == "text" {
		formatter = new(logrus.TextFormatter)
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	err = goerrors.Setup(formatter, level)
	if err != nil {
		panic(err)
	}
}
