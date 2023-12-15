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
		{
			Name:        "upload",
			Usage:       "upload a file to the remote implant",
			UsageText:   "upload [--origin] [--destination]",
			Description: "Running this command will read the file specified in [--origin] from the server and upload it to the [--destination] location on the remote implant.",
			Category:    "Modules",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "origin",
					Usage: "path to file on local server that should be uploaded",
				},
				&cli.StringFlag{
					Name:  "destination",
					Usage: "path on remote implant where file should be uploaded",
				},
			},
			Action: interactUpload,
		},
		{
			Name:        "download",
			Usage:       "download a file from the remote implant",
			UsageText:   "download [--origin] [--destination]",
			Description: "Running this command downloads a file from path [--origin] to the destination specified in [--destination].",
			Category:    "Modules",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "origin",
					Usage: "path to file on remote implant that should be downloaded",
				},
				&cli.StringFlag{
					Name:  "destination",
					Usage: "path on local server where file should be downloaded to",
				},
			},
			Action: interactDownload,
		},
		{
			Name:        "ls",
			Usage:       "list directory",
			UsageText:   "ls [path]",
			Description: "list the current working directory. If [path] is specified, lists the directory of [path].",
			Category:    "Modules",
			Action:      interactLs,
		},
		{
			Name:        "cd",
			Usage:       "change working directory",
			UsageText:   "cd [path]",
			Description: "changes the current working directory to [path].",
			Category:    "Modules",
			Action:      interactCd,
		},
		{
			Name:        "exec",
			Usage:       "execute a command",
			UsageText:   "exec [--cmd] [--args]",
			Description: "executes a command on the remote implant. This command is implemented using the \"os/exec\" package, see the docs for more information about the required parameters.",
			Category:    "Modules",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "cmd",
					Usage: "command to be executed",
				},
				&cli.StringFlag{
					Name:  "args",
					Usage: "arguments to program",
				},
			},
			Action: interactExec,
		},
		{
			Name:        "shellcode-exec",
			Usage:       "execute shellcode on remote implant",
			UsageText:   "shellcode-exec [shellcode]",
			Description: "executes shellcode on the remote implant. the shellcode should be in base64 format.",
			Category:    "Modules",
			Action:      interactShellcodeExec,
		},
	},
}
