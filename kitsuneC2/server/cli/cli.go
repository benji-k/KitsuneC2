//This package contains all CLI functionality. The CLI is responsible for parsing user input and executing server functionality
//through server/api.

package cli

import (
	"KitsuneC2/server/api"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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

func CliLoop() {
	for {
		if cliCtx.quit {
			return
		}

		if cliCtx.context == "home" {
			homeCliApp.Run(stringPrompt("KitsuneC2 > ", *color.New(color.FgRed)))
		} else if cliCtx.context == "interacting" {
			interactCliApp.Run(stringPrompt(cliCtx.implantId+" > ", *color.New(color.FgHiCyan)))
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
	r := bufio.NewReader(os.Stdin)
	for {
		c.Fprint(os.Stderr, label)
		userInput, _ = r.ReadString('\n')
		if userInput != "\n" { //The "\n" character is still in the buffer. This is basically the same as checking if the user provided input.
			break
		}
	}
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

//Template specific functions
//------------------homeCliApp functions-----------------------

func homeImplants(cCtx *cli.Context) error {
	implants, err := api.GetAllImplants()
	if err != nil {
		NotifyUser("could not fetch implants. Reason: "+err.Error(), "FAIL")
		return nil
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Implant ID", "IP", "Name", "Hostname", "User", "UID", "GID", "OS", "Arch", "Last Checkin")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for i := range implants {
		standardTimeFormat := time.Unix(int64(implants[i].Last_checkin), 0)
		tbl.AddRow(implants[i].Id, implants[i].Public_ip, implants[i].Name, implants[i].Hostname, implants[i].Username, implants[i].Uid, implants[i].Gid, implants[i].Os, implants[i].Arch, standardTimeFormat)
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
func interactModules(cCtx *cli.Context) error {

	return nil
}

func interactPendingTasks(cCtx *cli.Context) error {

	return nil
}

func interactCompletedTasks(cCtx *cli.Context) error {

	return nil
}

func interactAdd(cCtx *cli.Context) error {

	return nil
}

func interactRemove(cCtx *cli.Context) error {

	return nil
}

func interactKill(cCtx *cli.Context) error {

	return nil
}

func interactExit(cCtx *cli.Context) error {
	cliCtx.context = "home"
	cliCtx.implantId = ""

	return nil
}