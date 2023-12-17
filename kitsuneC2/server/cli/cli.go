//This package contains all CLI functionality. The CLI is responsible for parsing user input and executing server functionality
//through server/api.

package cli

import (
	"KitsuneC2/lib/communication"
	"KitsuneC2/lib/utils"
	"KitsuneC2/server/api"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/peterh/liner"
	"github.com/rodaine/table"
	cli "github.com/urfave/cli/v2"
)

// A cli can have 2 different contexts (can be expanced): a "home" page, or interacting with an implant.
type cliContext struct {
	context   string //either has value "home" or "interacting"
	implantId string //if context has value "interacting", this value will contains the implantId that is being interacted with
	quit      bool
}

var cliCtx cliContext = cliContext{context: "home", implantId: "", quit: false}
var rl *liner.State

func InitCli() {
	fmt.Println("Type 'help' for a list of available commands.")
	fmt.Println()
	rl = liner.NewLiner()
}

func CliLoop() {
	for {
		if cliCtx.quit {
			return
		}
		if cliCtx.context == "home" {
			homeCliApp.Run(stringPrompt("[KitsuneC2]", *color.New(color.FgRed)))
		} else if cliCtx.context == "interacting" {
			interactCliApp.Run(stringPrompt("["+cliCtx.implantId+"]", *color.New(color.FgHiCyan)))
		} else {
			log.Println("[ERROR] CLI is in unkown state, quitting!")
			return
		}
	}
}

// stringPrompt asks for a string value using the label in the color. Due to the way the urfave/cli package works,
// this function prepends the keyword "server" to a user provided string. This makes sure a user can immediately input
// sub-commands.
func stringPrompt(label string, c color.Color) []string {
	var result string = "server "
	var userInput string
	for {
		fmt.Println()
		c.Println(label)
		userInput, _ = rl.Prompt("> ")
		if userInput != "" {
			break
		}
	}
	rl.AppendHistory(userInput)
	result += userInput
	return strings.Split(strings.TrimSpace(result), " ")
}

// any error encountered in the CLI will get passed to this function.
func onCliError(cCtx *cli.Context, err error) {
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
	}
}

// Send a notification to the user in the CLI. The msgType can be one of the following values: "FAIL", "INFO", "SUCCESS".
func NotifyUser(msg string, msgType string) {
	switch msgType {
	case "FAIL":
		c := color.New(color.FgRed)
		c.Fprint(os.Stderr, "[*] ")
		fmt.Fprint(os.Stderr, msg+"\n")
	case "INFO":
		c := color.New(color.FgBlue)
		c.Fprint(os.Stdout, "[*] ")
		fmt.Fprint(os.Stdout, msg+"\n")
	case "SUCCESS":
		c := color.New(color.FgGreen)
		c.Fprint(os.Stdout, "[*] ")
		fmt.Fprint(os.Stdout, msg+"\n")
	default:
		c := color.New(color.FgBlue)
		c.Fprint(os.Stdout, "[*] ")
		fmt.Fprint(os.Stdout, msg+"\n")
	}
}

//Template specific functions. These functions get called from templates.go
//------------------homeCliApp functions-----------------------

func homeImplants(cCtx *cli.Context) error {
	implants, err := api.GetAllImplants()
	if err != nil {
		NotifyUser("could not fetch implants. Reason: "+err.Error(), "FAIL")
		return nil
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Implant ID", "IP", "Name", "Hostname", "User", "UID", "GID", "OS", "Arch", "Last Checkin (s)", "Active")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i := range implants {
		secondsSinceLastCheckin := time.Now().Unix() - implants[i].Last_checkin
		tbl.AddRow(implants[i].Id, implants[i].Public_ip, implants[i].Name, implants[i].Hostname, implants[i].Username, implants[i].Uid, implants[i].Gid, implants[i].Os, implants[i].Arch, secondsSinceLastCheckin, implants[i].Active)
	}
	tbl.Print()

	return nil
}

func homeGenerate(cCtx *cli.Context) error {

	return nil
}

func homeListenersList(cCtx *cli.Context) error {
	listeners, err := api.GetRunningListeners()
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Type", "Host", "Port")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i, listener := range *listeners {
		tbl.AddRow(i, listener.Type, listener.Network, listener.Port)
	}
	tbl.Print()
	return nil
}

func homeListenersAdd(cCtx *cli.Context) error {
	network := cCtx.String("host")
	port, err := strconv.Atoi(cCtx.String("port"))
	if err != nil {
		NotifyUser("port should be a valid integer", "FAIL")
		return nil
	}

	err = api.AddListener(network, port)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("listening on: "+network+":"+strconv.Itoa(port), "SUCCESS")
	return nil
}

func homeListenersRemove(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("remove: expected 1 argument", "FAIL")
		return nil
	}
	listenerId, err := strconv.Atoi(cCtx.Args().First())
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}

	err = api.KillListener(listenerId)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("Successfully stopped listener", "SUCCESS")
	return nil
}

func homeInteract(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("interact: expected 1 argument", "FAIL")
		return nil
	}
	userImplantId := cCtx.Args().First()
	if !api.ImplantExists(userImplantId) {
		NotifyUser("interact: no such implant", "FAIL")
		return nil
	}
	cliCtx.context = "interacting"
	cliCtx.implantId = userImplantId

	return nil
}

func homeQuit(cCtx *cli.Context) error {
	cliCtx.quit = true
	return nil
}

// ------------------interactCliApp functions-------------------

func interactPendingTasks(cCtx *cli.Context) error {
	tasks, err := api.GetTasksForImplant(cliCtx.implantId, false)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Module", "Arguments")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, task := range tasks {

		tbl.AddRow(task.Task_id, communication.MessageTypeToModuleName[task.Task_type], string(task.Task_data))
	}
	tbl.Print()
	return nil
}

func interactCompletedTasks(cCtx *cli.Context) error {
	tasks, err := api.GetTasksForImplant(cliCtx.implantId, true)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Module", "Arguments")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, task := range tasks {
		tbl.AddRow(task.Task_id, communication.MessageTypeToModuleName[task.Task_type], string(task.Task_data))
	}
	tbl.Print()
	return nil
}

func interactResult(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("remove: expected 1 argument", "FAIL")
		return nil
	}
	taskId := cCtx.Args().First()

	task, err := api.GetTask(taskId)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	if !task.Completed {
		NotifyUser("Task has no result because it hasn't been executed yet.", "FAIL")
		return nil
	}
	buff := bytes.NewBuffer([]byte{})
	json.Indent(buff, task.Task_result, "", "    ")
	fmt.Println(buff.String())
	return nil
}

func interactRemove(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("remove: expected 1 argument", "FAIL")
		return nil
	}
	taskId := cCtx.Args().First()
	err := api.RemovePendingTaskForImplant(cliCtx.implantId, taskId)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("Successfully removed task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactKill(cCtx *cli.Context) error {
	var task communication.Task = &communication.ImplantKillReq{}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 5, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactConfig(cCtx *cli.Context) error {
	serverIp := cCtx.String("server-ip")
	serverPort := cCtx.String("server-port")
	callbackInt := cCtx.String("callback-interval")
	callbackJit := cCtx.String("callback-jitter")

	var config = &communication.ImplantConfigReq{}
	if serverIp != "" {
		config.ServerIp = serverIp
	}
	serverPortI, err := strconv.Atoi(serverPort)
	if err != nil {
		config.ServerPort = -1
	} else {
		config.ServerPort = serverPortI
	}
	callbackIntI, err := strconv.Atoi(callbackInt)
	if err != nil {
		config.CallbackInterval = -1
	} else {
		config.CallbackInterval = callbackIntI
	}
	callbackJitI, err := strconv.Atoi(callbackJit)
	if err != nil {
		config.CallbackJitter = -1
	} else {
		config.CallbackJitter = callbackJitI
	}

	var task communication.Task = config
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 7, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")

	return nil
}

func interactExit(cCtx *cli.Context) error {
	cliCtx.context = "home"
	cliCtx.implantId = ""

	return nil
}

func interactFileInfo(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("remove: expected 1 argument", "FAIL")
		return nil
	}
	path := cCtx.Args().First()

	var task communication.Task = &communication.FileInfoReq{PathToFile: path}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 11, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactUpload(cCtx *cli.Context) error {
	origin := cCtx.String("origin")
	destination := cCtx.String("destination")
	if origin == "" || destination == "" {
		NotifyUser("both [--origin] and [--destination] must be valid paths.", "FAIL")
		return nil
	}

	fileContents, err := utils.ReadFile(origin)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}

	var task communication.Task = &communication.UploadReq{File: fileContents, Destination: destination}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 21, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactDownload(cCtx *cli.Context) error {
	origin := cCtx.String("origin")
	destination := cCtx.String("destination")
	if origin == "" || destination == "" {
		NotifyUser("both [--origin] and [--destination] must be valid paths.", "FAIL")
		return nil
	}

	var task communication.Task = &communication.DownloadReq{Origin: origin, Destination: destination}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 19, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactLs(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("ls: expected 1 argument", "FAIL")
		return nil
	}
	path := cCtx.Args().First()

	var task communication.Task = &communication.LsReq{Path: path}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 13, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactCd(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("cd: expected 1 argument", "FAIL")
		return nil
	}
	path := cCtx.Args().First()

	var task communication.Task = &communication.CdReq{Path: path}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 17, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")
	return nil
}

func interactExec(cCtx *cli.Context) error {
	cmd := cCtx.String("cmd")
	argsStr := cCtx.String("args")
	if cmd == "" {
		NotifyUser("[--cmd] is required", "FAIL")
		return nil
	}
	args := strings.Split(argsStr, " ")

	var task communication.Task = &communication.ExecReq{Cmd: cmd, Args: args}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 15, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")

	return nil
}

func interactShellcodeExec(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		NotifyUser("cd: expected 1 argument", "FAIL")
		return nil
	}
	scHex := cCtx.Args().First()
	sc, err := hex.DecodeString(scHex)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}

	var task communication.Task = &communication.ShellcodeExecReq{Shellcode: sc}
	taskId, err := api.AddTaskForImplant(cliCtx.implantId, 23, &task)
	if err != nil {
		NotifyUser(err.Error(), "FAIL")
		return nil
	}
	NotifyUser("created task with ID: "+taskId, "SUCCESS")

	return nil
}
