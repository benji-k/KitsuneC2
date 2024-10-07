//Package containing shared functionality for communication between an implant and the server. This file contains helper functions
//used by both the server and the implant to fascillitate networking functionality

package communication

import (
	"KitsuneC2/lib/cryptography"
	"encoding/json"
	"errors"
	"io"
	"net"
	"reflect"
)

// All communication between the client and servers get wrapped in an Envelope. This Envelope contains the type of the message being sent, and
// the message data itself. The Data variable can be further deserialized into specific message types.
type Envelope struct {
	MessageType int
	Data        []byte
}

// Given the type of message and it's arguments, wraps the data in an envelope datastructure and encrypts the whole datastructure. The function
// returns the encrypted representation of an envelope object in bytes.
func PackAndEncryptEnvelope(messageType int, data interface{}, aesKey []byte) ([]byte, error) {
	//First we check if the data variable corresponds to the correct messageType using reflection.
	expectedType := reflect.TypeOf(MessageTypeToStruct[messageType]())
	dataType := reflect.TypeOf(data)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return nil, errors.New("data does not correspond to messageType")
	}

	//After we know all types are correct, marshal the passed data and put it in an envelope. Afterwards marshal the whole envelope so that
	//it can be encrypted.
	envelopeData, _ := json.Marshal(data)
	rawJson, _ := json.Marshal(Envelope{MessageType: messageType, Data: envelopeData})

	encryptedJson, err := cryptography.EncryptAes(rawJson, []byte(aesKey))
	if err != nil {
		return nil, err
	}
	return encryptedJson, nil
}

// Writes data to a socket. The function checks if the number of bytes actually written to the socket == len(data).
func WriteToSocket(connection net.Conn, data []byte) error {
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
func ReadFromSocket(connection net.Conn, n int) ([]byte, error) {
	var buffer []byte = make([]byte, n)
	bytesRead, err := io.ReadFull(connection, buffer)
	if err != nil {
		return nil, err
	}
	if bytesRead != n {
		return nil, errors.New("data was only partially received from server")
	}
	return buffer, nil
}
