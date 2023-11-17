//Package containing functionality for communication between an implant and the server. This file contains all structures that are sent over the wire
//by the implant/server. These structs get serialized/deserialized with the JSON package.

package communication

// This map correlates all MessageType's to a specific data stucture for a message. This can be used for reflection so that no big switch
// statements have to be created
var MessageTypeToStruct = map[int]func() interface{}{
	0: func() interface{} { return &ImplantRegister{} },
	1: func() interface{} { return &ImplantCheckin{} },
	//reserved for implant functionality
	11: func() interface{} { return &FileInfoReq{} },
	12: func() interface{} { return &FileInfoResp{} },
}

type ImplantRegister struct {
	ImplantId   string
	ImplantName string
	Hostname    string
	Username    string
	UID         string
	GID         string
}

type ImplantCheckin struct {
	ImplantId string
}

type FileInfoReq struct {
	PathToFile string
}

type FileInfoResp struct {
	Name    string
	Size    int64
	Mode    string
	ModTime int
	IsDir   bool
}
