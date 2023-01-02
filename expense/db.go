package expense

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDB() {
	// os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", "postgres://rudwxbbq:UX7e7F375OJMZDzvtMd5BWlwenaRM0mv@tiny.db.elephantsql.com/rudwxbbq")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	command := "CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);"

	_, err = db.Exec(command)
	if err != nil {
		log.Fatal("Unable to create table", err)
	}
}
