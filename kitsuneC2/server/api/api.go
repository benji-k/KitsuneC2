package api

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/cryptography"
	"KitsuneC2/server/db"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

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
