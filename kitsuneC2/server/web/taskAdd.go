package web

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/server/api"

	"github.com/gin-gonic/gin"
)

func addImplantKill(c *gin.Context, implantId string) error {
	var task communication.Task = &communication.ImplantKillReq{}
	_, err := api.AddTaskForImplant(implantId, 5, &task)
	if err != nil {
		return err
	}
	return nil
}

func addChangeConfig(c *gin.Context, implantId string) error {

	return nil
}

func addFileInfo(c *gin.Context, implantId string) error {

	return nil
}

func addLs(c *gin.Context, implantId string) error {

	return nil
}

func addExec(c *gin.Context, implantId string) error {

	return nil
}

func addCd(c *gin.Context, implantId string) error {

	return nil
}

func addDownload(c *gin.Context, implantId string) error {

	return nil
}

func addUpload(c *gin.Context, implantId string) error {

	return nil
}

func addShellcodeExec(c *gin.Context, implantId string) error {

	return nil
}
