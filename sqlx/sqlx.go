package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

const (
	dbHost = "DBHOST"
	dbPort = "DBPORT"
	dbUser = "DBUSER"
	dbPwd  = "DBPASS"
	dbName = "DBNAME"
)

// Todo list ...
type Todo struct {
	TaskID      int          `db:"taskid"`
	Topic       string       `db:"topic"`
	Description string       `db:"description"`
	Completed   bool         `db:"completed"`
	CreatedAt   time.Time    `db:"createdat"`
	UpdatedAt   sql.NullTime `db:"updatedat"`
}

// // NullTime is an alias for sql.NullTime data type
// type NullTime struct {
// 	sql.NullTime
// }

// // MarshalSON for NullTime
// func (nt *NullTime) MarshalSON() ([]byte, error) {
// 	if !nt.Valid {
// 		return []byte("null"), nil
// 	}
// 	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
// 	return []byte(val), nil
// }

func main() {

	initDB()
	defer db.Close()

	// Simulating Atomicity when inserting records into the database.
	// transactions begins
	tx := db.MustBegin()
	now := time.Now()
	t := Todo{
		Topic:       "Rock Star Quest",
		Description: "Learn the electric guitar!",
		Completed:   false,
		CreatedAt:   now,
	}
	tx.Exec(`insert into "TODOTestSchema".todos (topic, description, completed, createdAt) values ($1, $2, $3, $4)`,
		t.Topic, t.Description, t.Completed, t.CreatedAt)
	tx.Commit()
	// transaction ends

	// dealing with NULLable columns
	myTodos := []Todo{}
	db.Select(&myTodos, `select * from "TODOTestSchema".todos`)

	for _, todo := range myTodos {

		log.Printf("TaskID: %v\n", todo.TaskID)
		log.Printf("Topic: %v\n", todo.Topic)
		log.Printf("Description: %v\n", todo.Description)
		log.Printf("Completed: %v\n", todo.Completed)
		log.Printf("CreatedAt: %v\n", todo.CreatedAt)
		log.Printf("UpdatedAt: %v\n", todo.UpdatedAt)
	}

}

func initDB() {
	config := dbConfig()

	var err error

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbHost], config[dbPort], config[dbUser], config[dbPwd], config[dbName])

	db, err = sqlx.Open("postgres", psqlConn)

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
