//Package containing functionality for communication between an implant and the server. This file contains all structures that are sent over the wire
//by the implant/server. These structs get serialized/deserialized with the JSON package.

package communication

// This map correlates all MessageType's to a specific data stucture for a message. This can be used for reflection so that no big switch
// statements have to be created. Note that MessageTypes and TaskTypes are the same.
var MessageTypeToStruct = map[int]func() interface{}{
	0: func() interface{} { return &ImplantRegister{} },
	1: func() interface{} { return &ImplantCheckinReq{} },
	2: func() interface{} { return &ImplantCheckinResp{} },
	4: func() interface{} { return &ImplantErrorResp{} },
	//reserved for implant functionality
	11: func() interface{} { return &FileInfoReq{} },
	12: func() interface{} { return &FileInfoResp{} },
}

// A task is a type of message that gets sent after checkin. Every task needs to have an ID so that we know to what task certain
// output belongs.
type Task interface {
	SetTaskId(id string) //This function is only used in the API. It's the API responsibility to generate tasks.
}

type ImplantRegister struct {
	ImplantId   string
	ImplantName string
	Hostname    string
	Username    string
	UID         string
	GID         string
}

type ImplantCheckinReq struct {
	ImplantId string
}

type ImplantCheckinResp struct {
	TaskTypes     []int
	TaskArguments [][]byte
}

type ImplantErrorResp struct {
	TaskId string
	Error  string
}

type FileInfoReq struct {
	TaskId     string
	PathToFile string
}

func (t *FileInfoReq) SetTaskId(id string) {
	t.TaskId = id
}

type FileInfoResp struct {
	TaskId  string
	Name    string
	Size    int64
	Mode    string
	ModTime int
	IsDir   bool
}

func (t *FileInfoResp) SetTaskId(id string) {
	t.TaskId = id
}
