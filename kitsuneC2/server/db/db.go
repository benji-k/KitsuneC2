package db

import (
	//"github.com/mattn/go-sqlite3"
	"database/sql"
)

var (
	dbConn *sql.DB
)

func Init() {
	dbConn, _ = sql.Open("sqlite3", "CHANGEME")
}

func Shutdown() {

}

func AddImplant() {

}

func RemoveImplant() {

}

func AddListener() {

}

func RemoveListener() {

}

func AddPayload() {

}

func RemovePayload() {

}

func AddTask() {

}

func RemoveTask() {

}
