package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Attempts to read $HOME/.kitsuneC2/kitsuneC2.log and creates the file if it doesn't exist. If this fails, logs will be written
// to stderr
func InitLogger() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("[ERROR] logging: could not find user home directory to write log files to. Logs will be discarded.")
		log.SetOutput(io.Discard)
		return
	}
	_, err = os.Stat(filepath.Join(homeDir, ".kitsuneC2"))
	if err != nil {
		os.Mkdir(filepath.Join(homeDir, ".kitsuneC2"), 0700)
	}

	var logFile *os.File

	_, err = os.Stat(filepath.Join(homeDir, ".kistuneC2", "kitsuneC2.log"))
	if err != nil {
		logFile, err = os.Create(filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
		if err != nil {
			log.Printf("[ERROR] logging: could not write logs to %s. Logs will be discarded", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
			log.SetOutput(io.Discard)
			return
		}
	} else {
		logFile, err = os.Open(filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
		if err != nil {
			log.Printf("[ERROR] logging: could not write logs to %s. Logs will be discarded", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
			log.SetOutput(io.Discard)
			return
		}
	}

	log.Printf("[INFO] logging: Logging to %s.", filepath.Join(homeDir, ".kitsuneC2", "kitsuneC2.log"))
	log.SetOutput(logFile)
}
