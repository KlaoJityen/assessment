package expense

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func InitDB() {

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	command := "CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);"

	_, err = db.Exec(command)
	if err != nil {
		log.Fatal("Unable to create table", err)
	}
}
