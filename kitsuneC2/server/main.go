package main

import (
	"KitsuneC2/server/db"
	"KitsuneC2/server/listener"
	"time"
)

func main() {
	db.Init()
	l1 := listener.Listener{Type: "tcp", Handler: tcpHandler, Network: "127.0.0.1", Port: 4444}
	l1.Start()
	time.Sleep(time.Minute * 10)
	db.Shutdown()
}
