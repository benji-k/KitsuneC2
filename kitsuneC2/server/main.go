package main

import (
	"KitsuneC2/lib/utils"
	"KitsuneC2/server/cli"
	"KitsuneC2/server/db"
	"KitsuneC2/server/logging"
	"KitsuneC2/server/transport"
	"KitsuneC2/server/web"
	"os"

	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() //loads variables in .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initialize()
	cli.CliLoop()
	db.Shutdown()
}

func initialize() {
	utils.PrintBanner()
	db.Initialize()
	cli.InitCli()
	transport.Initialize()
	logging.InitLogger()

	if os.Getenv("ENABLE_WEB_API") == "true" {
		web.Init()
	}
}
