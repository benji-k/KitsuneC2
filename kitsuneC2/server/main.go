package main

import (
	"os"

	"github.com/benji-k/KitsuneC2/kitsuneC2/lib/utils"
	"github.com/benji-k/KitsuneC2/kitsuneC2/server/cli"
	"github.com/benji-k/KitsuneC2/kitsuneC2/server/db"
	"github.com/benji-k/KitsuneC2/kitsuneC2/server/logging"
	"github.com/benji-k/KitsuneC2/kitsuneC2/server/transport"
	"github.com/benji-k/KitsuneC2/kitsuneC2/server/web"

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
