//This package serves as the main API to the database. Things like implant info, pending tasks etc. can be found in this package.
//This is also the only package that contains raw SQL commands. The package is structured as follows:
//db.go: Functionality related to creating the DB and opening a connection to it
//application.go: Contains all API functions that other parts of the program make calls to
//structs.go: Calling functions must call the API functions with structs defined in this file.
//sql.schema: When no database file exists, the package attempts to create one based on this schema file.

package db

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB
var dbPath string

// Embed the schema.sql file
//
//go:embed sql.schema
var schema string

// Initializes the database by looking for the kitsune.sqlite specified in dbPath variable. If it finds it, it opens a connection to it.
// If it doesn't, it attempts to create a database file. This function must succeed for the program to function properly.
func Initialize() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("[FATAL] db: user home directory not available to write files to!")
	}
	dbPath = filepath.Join(homeDir, ".kitsuneC2", "kitsune.sqlite")

	//Check if db already exists. If not, create it.
	if _, err := os.Stat(dbPath); err != nil {
		log.Printf("[INFO] db: %s was not found, attempting to create...", dbPath)
		err := createDb()
		if err != nil {
			log.Fatal("[FATAL] db: could not create database file! Reason: " + err.Error())
		}
		log.Printf("[INFO] db: succesfully created %s", dbPath)
	}

	log.Printf("[INFO] db: %s found, opening connection...", dbPath)

	dbConn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("[FATAL] db: could not connect to db! Reason: " + err.Error())
	}
	log.Println("[INFO] db: Database initialized and ready to serve...")
}

func Shutdown() {
	log.Println("[INFO] db: Shutting down DB...")
	dbConn.Close()
}

// When no kitsune.sqlite file is found, this function attempts to create one based on the sql.schema file.
func createDb() error {

	// Open a connection to the SQLite database (this will create the file if it doesn't exist)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the schema
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
