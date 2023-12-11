package main

import (
	"KitsuneC2/server/db"
	//"KitsuneC2/server/listener"
	//"time"

	"KitsuneC2/server/cli"
)

func main() {
	db.Init()
	cli.CliLoop()
	db.Shutdown()
}
