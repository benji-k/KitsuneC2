package transport

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"KitsuneC2/server/db"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
)

// A Session consists of an open socket connection, and the aes-key used to encrypt/decrypt data going in/out
// this socket
type Session struct {
	Connection net.Conn
	AesKey     []byte
}

var privKey string

// Fetches (if exists) keypair from database for implant communication. If no keypair exists, attempts to create one.
func Initialize() {
	priv, err := db.GetPrivateKey()
	if err == db.ErrNoResults { //we don't have a keypair, try to create one.
		priv, pub, err := cryptography.GenerateRSAKeyPair(2048)
		if err != nil {
			log.Fatal("Could not generate keypair. Reason: ", err.Error())
		}
		privStr := cryptography.RsaPrivateKeyToString(priv)
		pubStr := cryptography.RSAPublicKeyToString(pub)
		db.InitKeypair(privStr, pubStr)
	} else if err != nil { //something went wrong during fetching of keypair
		log.Fatal("Could not fetch private key from db. Reason: " + err.Error())
	}
	//we managed to fetch private key, or create a keypair

	privKey = priv
}

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
