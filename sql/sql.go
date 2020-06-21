package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbHost = "DBHOST"
	dbPort = "DBPORT"
	dbUser = "DBUSER"
	dbPwd  = "DBPASS"
	dbName = "DBNAME"
)

func main() {

	initDB()
	defer db.Close()

	var todoTitle string

	// rows, err := db.Query("select title from todos")
	rows, err := db.Query(`select title from "TODOTestSchema".todos`)

	for rows.Next() {
		rows.Scan(&todoTitle)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("TODO: %s\n", todoTitle)
	}

}

func initDB() {
	config := dbConfig()

	var err error

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbHost], config[dbPort], config[dbUser], config[dbPwd], config[dbName])

	db, err = sql.Open("postgres", psqlConn)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")

}

func dbConfig() map[string]string {

	config := make(map[string]string)

	host, found := os.LookupEnv(dbHost)
	if !found {
		panic("DBHOST environment variable not provided.")
	}
	port, found := os.LookupEnv(dbPort)
	if !found {
		panic("DBPORT environment variable not provided.")
	}
	user, found := os.LookupEnv(dbUser)
	if !found {
		panic("DBUSER environment variable not provided.")
	}
	pwd, found := os.LookupEnv(dbPwd)
	if !found {
		panic("DBPASS environment variable not provided.")
	}
	name, found := os.LookupEnv(dbName)
	if !found {
		panic("DBNAME environment variable not provided.")
	}

	config[dbHost] = host
	config[dbPort] = port
	config[dbUser] = user
	config[dbPwd] = pwd
	config[dbName] = name

	return config

}
