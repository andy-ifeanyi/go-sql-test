package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

// Todo list ...
type Todo struct {
	TaskID      int
	Topic       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	// UpdatedAt   time.Time
}

func main() {

	initDB()
	defer db.Close()

	// var todoTitle string

	// query database for a single column
	// rows, err := db.Query(`select title from "TODOTestSchema".todos`) // rows, err := db.Query("select title from todos")

	// query database for multiple columns
	rows, err := db.Query(`select taskid, topic, description, completed, createdat from "TODOTestSchema".todos`)
	for rows.Next() {

		todo := Todo{}

		if err = rows.Scan(&todo.TaskID, &todo.Topic, &todo.Description, &todo.Completed, &todo.CreatedAt); err != nil {
			log.Fatal(err)
		}

		log.Printf("TaskID: %v\n", todo.TaskID)
		log.Printf("Topic: %v\n", todo.Topic)
		log.Printf("Description: %v\n", todo.Description)
		log.Printf("Completed: %v\n", todo.Completed)
		log.Printf("CreatedAt: %v\n", todo.CreatedAt)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// insert into database
	now := time.Now() // used to create createdAt / updatedAt timestamps

	res, err := db.Exec(`insert into "TODOTestSchema".todos (topic, description, completed, createdAt, updatedAt) values ($1, $2, $3, $4, $5)`,
		"Multivariate Calcalus", "Hone up my Machine Learning Kung Fu", false, now, now)
	if err != nil {
		log.Fatal(err)
	}
	affected, _ := res.RowsAffected()
	log.Printf("Rows affected %d\n", affected)

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
