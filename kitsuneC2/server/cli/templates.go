package cli

//This file contains all the different cli "apps" (states) that the CLI can be in. If we are currently interacting with an implant,
//there must be different commands available then when you just started the server.

import (
	cli "github.com/urfave/cli/v2"
)

// This app gets executed when the user is on the 'homeCliApp' page, e.g. when the user just starts the server.
var homeCliApp cli.App = cli.App{
	Name:           " ",
	Usage:          " ",
	UsageText:      "[command] [sub-command] [arguments...]",
	Description:    "KitsuneC2 server. See list below for available commands. For more information about a command, type \"[command] help\".",
	ExitErrHandler: onCliError,
	Commands: []*cli.Command{
		{
			Name:        "implants",
			Usage:       "List active implants",
			UsageText:   "implants",
			Description: "Lists all active implants that have contacted the server.",
			Action:      homeImplants,
		},
		{
			Name:        "generate",
			Usage:       "Generate a new KistuneC2 implant binary",
			UsageText:   "generate [arguments]",
			Description: "Generates a new KitsuneC2 implant binary.",
			Action:      homeGenerate,
		},
		{
			Name:        "listeners",
			Usage:       "Add or remove listeners",
			UsageText:   "listeners [command]",
			Description: "Add or remove a TCP listener",
			Subcommands: []*cli.Command{
				{
					Name:      "list",
					Usage:     "list all running listeners",
					UsageText: "list",
					Action:    homeListenersList,
				},
				{
					Name:      "add",
					Usage:     "add a new listener",
					UsageText: "listeners add [--host] [--port]",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "host",
							Value: "0.0.0.0",
							Usage: "interface on which server will be listening for connections.",
						},
						&cli.StringFlag{
							Name:  "port",
							Value: "4444",
							Usage: "port on which server will be listening for connections.",
						},
					},
					Action: homeListenersAdd,
				},
				{
					Name:      "remove",
					Usage:     "remove an existing listener",
					UsageText: "remove [listener ID]",
					Action:    homeListenersRemove,
				},
			},
		},
		{
			Name:        "interact",
			Usage:       "Interact with specific implant",
			UsageText:   "interact [implant ID]",
			Description: "List active tasks, see results of executed tasks, and control implant with [implant ID]",
			Action:      homeInteract,
		},
		{
			Name:        "quit",
			Usage:       "Quits KitsuneC2",
			UsageText:   "quit",
			Description: "Quits KitsuneC2",
			Action:      homeQuit,
		},
	},
}

// this app gets executed when a user is interacting with a specific implant
var interactCliApp cli.App = cli.App{
	Name:           " ",
	Usage:          " ",
	UsageText:      "[command] [sub-command] [arguments...]",
	Description:    "KitsuneC2 implant. See list below for available commands. For more information about a command, type \"[command] help\".",
	ExitErrHandler: onCliError,
	Commands: []*cli.Command{
		{
			Name:        "pending-tasks",
			Usage:       "List all pending tasks for this implant",
			UsageText:   "pending-tasks",
			Description: "Lists all pending tasks for this implant.",
			Action:      interactPendingTasks,
		},
		{
			Name:        "completed-tasks",
			Usage:       "List all completed tasks for this implant",
			UsageText:   "completed-tasks",
			Description: "Lists all completed tasks for this implant.",
			Action:      interactCompletedTasks,
		},
		{
			Name:        "result",
			Usage:       "Check the output of an executed task",
			UsageText:   "result [task ID]",
			Description: "Returns the result for an executed task.",
			Action:      interactResult,
		},
		{
			Name:        "remove",
			Usage:       "Remove task from the list of pending tasks",
			UsageText:   "remove [task ID]",
			Description: "Removes a task from the list of pending tasks. The implant will not execute this task anymore on the next check-in.",
			Action:      interactRemove,
		},
		{
			Name:        "kill",
			Usage:       "Kill this implant",
			UsageText:   "kill",
			Description: "Running this command will delete the implant from the infected host.",
			Action:      interactKill,
		},
		{
			Name:        "exit",
			Usage:       "Stop interacting with this implant",
			UsageText:   "exit",
			Description: "Stops interaction with this implant and returns to home screen.",
			Action:      interactExit,
		},
		{
			Name:        "file-info",
			Usage:       "get information about a file",
			UsageText:   "file-info [path]",
			Description: "Running this command will fetch information about a file on the remote host.",
			Category:    "Modules",
			Action:      interactFileInfo,
		},
	},
}
