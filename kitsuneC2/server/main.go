package main

import (
	"KitsuneC2/lib/utils"
	"KitsuneC2/server/cli"
	"KitsuneC2/server/db"
	"KitsuneC2/server/transport"
)

func main() {
	initialize()
	cli.CliLoop()
	db.Shutdown()
}

func initialize() {
	utils.PrintBanner()
	db.Initialize()
	cli.InitCli()
	transport.Initialize()
}
