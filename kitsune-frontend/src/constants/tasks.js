export const Tasks = [
    {
        taskType: 13,
        taskName: "ls",
        description: "Lists directory contents.",
        args: [
            {
                name: "path",
                apiName: "Path",
                optional: false,
                type: String,
                tooltip: "If not specified, lists current working directory."
            }
        ]
    },
    {
        taskType: 5,
        taskName: "implant kill",
        description: "Kills the implant. [DANGEROUS]",
        args: []
    },
    {
        taskType: 7,
        taskName: "change config",
        description: "Change the implant configuration.",
        args: [
            {
                name: "server-ip",
                apiName: "ServerIp",
                optional: true,
                type: String,
                tooltip: "New server IP that implant will connect to."
            },
            {
                name: "server-port",
                apiName: "ServerPort",
                optional: true,
                type: Number,
                tooltip: "New server port that implant will connect to."
            },
            {
                name: "callback-interval",
                apiName: "CallbackInterval",
                optional: true,
                type: Number,
                tooltip: "Time between implant check-ins."
            },
            {
                name: "callback-jitter",
                apiName: "CallbackJitter",
                optional: true,
                type: Number,
                tooltip: "Randomness of implant checkins in seconds."
            },

        ]
    },
    {
        taskType: 11,
        taskName: "file info",
        description: "Fetch info, such as file-size, permissions etc. from a specific file.",
        args: [
            {
                name: "path",
                apiName: "PathToFile",
                optional: false,
                type: String,
                tooltip: "Path to file you want more info of."
            }
        ]
    },
    {
        taskType: 15,
        taskName: "exec",
        description: "Executes a program with specified arguments on the remote machine.",
        args: [
            {
                name: "cmd",
                apiName: "Cmd",
                optional: false,
                type: String,
                tooltip: "Program to be executed, e.g. bash. (See Go's 'exec' documentation for details)"
            },
            {
                name: "args",
                apiName: "Args",
                optional: true,
                type: String,
                tooltip: "Arguments to be passed to executing program."
            }
        ]
    },
    {
        taskType: 17,
        taskName: "cd",
        description: "Changes the working directory of the implant.",
        args: [
            {
                name: "path",
                apiName: "Path",
                optional: false,
                type: String,
                tooltip: "Path to new working directory."
            }
        ]
    },
    {
        taskType: 19,
        taskName: "download",
        description: "Downloads a file from the remote host.",
        args: [
            {
                name: "path",
                apiName: "Origin",
                optional: false,
                type: String,
                tooltip: "Path to file that should be downloaded."
            }
        ]
    },
    {
        taskType: 21,
        taskName: "upload",
        description: "Uploads a file to the remote host.",
        args: [
            {
                name: "file",
                apiName: "File",
                optional: false,
                type: "file",
                tooltip: "File that should be uploaded"
            },
            {
                name: "destination",
                apiName: "Destination",
                optional: false,
                type: String,
                tooltip: "Directory file should be uploaded to."
            }
        ]
    },
    {
        taskType: 23,
        taskName: "shellcode exec",
        description: "Executes shellcode in a new thread on the remote host.",
        args: [
            {
                name: "shellcode",
                apiName: "Shellcode",
                optional: false,
                type: String,
                tooltip: "Base64 encoded shellcode"
            }
        ]
    },
]