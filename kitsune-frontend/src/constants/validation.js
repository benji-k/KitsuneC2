//All API endpoints that receive parameters from the client, have them validated using validate.js. The requirements of the parameters are defined here

export const Validation = {
    api_kitsune_implants_generate : {
        os : {
            presence: {allowEmpty: false},
            inclusion: { within: ["linux", "windows", "aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios", "js", "netbsd", "plan9", "solaris"] }
        },
        arch : {
            presence: {allowEmpty: false},
            inclusion: { within: ["386", "amd64", "arm", "arm64", "mips", "mips64", "mips64le", "mipsle", "ppc64", "ppc64le", "riscv64", "s390x", "wasm"] }
        },
        serverIp : {
            presence: {allowEmpty: false},
            format: {
                pattern: /^(?:\d{1,3}\.){3}\d{1,3}$|^[a-fA-F0-9:]+$/,  // IPv4 or IPv6 pattern
                message: "must be a valid IP address (either IPv4 or IPv6)",
            },
        },
        name : {    
 
        },
        serverPort : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
                lessThanOrEqualTo: 65535,
            }
        },
        cbInterval : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
            }
        },
        cbJitter : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
            }
        },
        maxRetryCount : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
            }
        }
    },

    api_kitsune_implants_remove : {
        implants : {
            presence: {allowEmpty: false},
        }
    },

    api_kitsune_listeners_add : {
        network : {
            presence: {allowEmpty: false},
            format: {
                pattern: /^(?:\d{1,3}\.){3}\d{1,3}$|^[a-fA-F0-9:]+$/,  // IPv4 or IPv6 pattern
                message: "must be a valid IP address (either IPv4 or IPv6)",
            },
        },
        port : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
                lessThanOrEqualTo: 65535,
            }
        }
    },

    api_kitsune_listeners_remove : {
        id : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThanOrEqualTo: 0,
            }
        }
    },

    api_kitsune_tasks_add : {
        taskType : {
            presence: {allowEmpty: false},
            numericality: {
                onlyInteger: true,
                greaterThan: 0,
            },
            inclusion: { within: ["5", "7", "11", "13", "15", "17", "19", "21", "23"]} //see constants/tasks.js for all available taskTypes
        },
        implants : {
            presence: {allowEmpty: false},
        }
    },

    api_kitsune_tasks_remove : {
        implantId : {
            presence: {allowEmpty: false},
            format: {
                pattern: "^[a-zA-Z0-9]*$",  // Alphanumeric (letters and numbers) or empty
                message: "can only contain letters and numbers",
            },
        },
        taskId:{
            presence: {allowEmpty: false},
            format: {
                pattern: "^[a-zA-Z0-9]*$",  // Alphanumeric (letters and numbers) or empty
                message: "can only contain letters and numbers",
            },
        }
    },

    api_kitsune_file_download:{
        taskId:{
            presence: {allowEmpty: false},
            format: {
                pattern: "^[a-zA-Z0-9]*$",  // Alphanumeric (letters and numbers) or empty
                message: "can only contain letters and numbers",
            },
        }
    }
}