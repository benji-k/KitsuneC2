package logging

import (
	"log"
	"os"
	"path/filepath"
)

var logFile *os.File

// Attempts to read $HOME/.kitsuneC2/kitsuneC2.log and creates the file if it doesn't exist. If this fails, logs will be written
// to stderr
func InitLogger() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[ERROR] logging: could not find user home directory to write log files to. Logs will be discarded.")
		log.SetOutput(os.Stderr)
		return
	}
	_, err = os.Stat(filepath.Join(homeDir, ".kitsuneC2"))
	if err != nil {
		os.Mkdir(filepath.Join(homeDir, ".kitsuneC2"), 0700)
	}

	_, err = os.Stat(filepath.Join(homeDir, ".kistuneC2", "kitsuneC2.log"))
	if err != nil {
		logFile, err = os.Create(filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
		if err != nil {
			log.Printf("[ERROR] logging: could not write logs to %s. Logs will be written to stderr", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
			log.SetOutput(os.Stderr)
			return
		}
	} else {
		logFile, err = os.OpenFile(filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"), os.O_APPEND, 0600)
		if err != nil {
			log.Printf("[ERROR] logging: could not write logs to %s. Logs will be written to stderr", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
			log.SetOutput(os.Stderr)
			return
		}
	}

	log.Printf("[INFO] logging: Logging to %s.", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
	log.SetOutput(logFile)
}

func ShutdownLogger() {
	logFile.Close()
}

func GetLogFilepath() string {
	return logFile.Name()
}
