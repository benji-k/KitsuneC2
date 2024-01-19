package web

import (
	"KitsuneC2/server/api"
	"KitsuneC2/server/db"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Init() {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	gin.LoggerWithWriter(log.Writer())

	router.GET("/implants", getImplants)
	router.GET("/tasks", getTasks)
	go router.Run("0.0.0.0:7331")
}

func getImplants(c *gin.Context) {
	implants, err := api.GetAllImplants()
	if err == db.ErrNoResults {
		emptyResp := make([]string, 0)
		c.JSON(200, emptyResp)
		return
	} else if err != nil {
		returnError(c, err)
		return
	}
	c.JSON(200, implants)
}

func getTasks(c *gin.Context) {
	completed := false
	paramAsStr := c.Request.URL.Query().Get("completed")
	if paramAsStr != "" {
		paramAsBool, err := strconv.ParseBool(paramAsStr)
		if err == nil {
			completed = paramAsBool
		}
	}
	tasks, err := api.GetAllTasks(completed)
	if err == db.ErrNoResults {
		emptyResp := make([]string, 0)
		c.JSON(200, emptyResp)
		return
	} else if err != nil {
		returnError(c, err)
		return
	}
	c.JSON(200, tasks)
}

func returnError(c *gin.Context, err error) {
	c.JSON(500, gin.H{"error": err.Error()})
}
