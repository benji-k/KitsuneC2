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
)

// Sends an "envelope" (see /lib/communication) to the server. This function encryptes the envelope with AES using the session key.
// Before sending the data, the function prepends it with the length of the data, so the server knows how much data to receive.
func SendEnvelope(connection net.Conn, envelope Envelope, aesKey []byte) error {
	rawJson, _ := json.Marshal(envelope)
	encryptedJson, err := cryptography.EncryptAes(rawJson, []byte(aesKey))
	if err != nil {
		return err
	}

	//create a buffer and write the len(encryptedData) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)))
	buffer.Write(encryptedJson)

	err = WriteToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Reads the first 4 bytes (uint32) from socket (message length). The reads the entire encrypted message. Afterwards, the function decrypts the message
// and attempts to unmarshal the resulting JSON into a Envelope object.
func ReceiveEnvelope(connection net.Conn, aesKey []byte) (*Envelope, error) {
	messageLengthAsBytes, err := ReadFromSocket(connection, 4)
	if err != nil {
		return nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	cipherText, err := ReadFromSocket(connection, int(messageLength))
	if err != nil {
		return nil, err
	}
	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, []byte(aesKey))
	if err != nil {
		return nil, err
	}

	var messageEnvelope *Envelope = new(Envelope)
	json.Unmarshal(rawJsonAsBytes, messageEnvelope)

	return messageEnvelope, nil
}

// writes data to a socket. The function checks if the number of bytes actually written to the socket == len(data).
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

// reads n bytes from a socket.
func ReadFromSocket(connection net.Conn, n int) ([]byte, error) {
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
