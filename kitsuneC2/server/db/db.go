//This package serves as the main API to the database. Things like implant info, pending tasks etc. can be found in this package.
//This is also the only package that contains raw SQL commands. The package is structured as follows:
//db.go: Functionality related to creating the DB and opening a connection to it
//application.go: Contains all API functions that other parts of the program make calls to
//structs.go: Calling functions must call the API functions with structs defined in this file.
//sql.schema: When no database file exists, the package attempts to create one based on this schema file.

package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dbConn *sql.DB
)

// Initializes the database by looking for a kitsune.sqlite file in the /db/ folder. If it finds it, it opens a connection to it.
// If it doesn't, it attempts to create a database file. This function must succeed for the program to function properly.
func Init() {
	//Check if db already exists. If not, create it.
	if _, err := os.Stat("./db/kitsune.sqlite"); err != nil {
		log.Println("[INFO] /db/kitsune.sqlite was not found, attempting to create...")
		err := createDb()
		if err != nil {
			log.Fatal("[FATAL] could not create database file! Reason: " + err.Error())
		}
		log.Println("[INFO] succesfully created /db/kitsune.sqlite")
	}

	log.Println("[INFO] /db/kitsune.sqlite found, opening connection...")

	var err error
	dbConn, err = sql.Open("sqlite3", "./db/kitsune.sqlite")
	if err != nil {
		log.Fatal("[FATAL] could not connect to db! Reason: " + err.Error())
	}
	log.Println("[INFO] Database initialized and ready to serve...")
}

func Shutdown() {
	dbConn.Close()
}

// When no kitsune.sqlite file is found, this function attempts to create one based on the sql.schema file.
func createDb() error {
	schema, err := os.Open("./db/sql.schema")
	if err != nil {
		return err
	}
	defer schema.Close()

	cmd := exec.Command("sqlite3", "./db/kitsune.sqlite")
	cmd.Stdin = schema

	cmdOutput, err := cmd.Output()
	if string(cmdOutput) != "" || err != nil {
		os.Remove("./db/kitsune.sqlite")
		return errors.New("\"sqlite3 kitsude.sqlite < sql.schema\" gave an unexpected output! Output: " + string(cmdOutput) + ". Error: " + err.Error())
	}
	return nil
}
