package main

import (
	"KitsuneC2/lib/utils"
	"KitsuneC2/server/cli"
	"KitsuneC2/server/db"
	"KitsuneC2/server/logging"
	"KitsuneC2/server/transport"
	"KitsuneC2/server/web"
	"os"

	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENABLE_WEB_API") == "true" {
		fmt.Println("Starting in daemon mode. Type \"exit\" to shutdown daemon.")
		err := godotenv.Load() //loads variables in .env file
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		initialize()
		cmd := ""
		for cmd != "exit" {
			fmt.Scanln(&cmd)
		}
	} else {
		initialize()
		cli.CliLoop()
	}

	shutdown()
}

func initialize() {
	utils.PrintBanner()
	logging.InitLogger()
	db.Initialize()
	if os.Getenv("ENABLE_WEB_API") == "true" {
		web.Init()
	} else {
		cli.InitCli()
	}
	transport.Initialize()
}

func shutdown() {
	db.Shutdown()
	logging.ShutdownLogger()
}
