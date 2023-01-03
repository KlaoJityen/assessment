package expense

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var handler *Handler

func InitDB() {

	var err error
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	handler = NewApplication(db)

	command := "CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);"

	_, err = handler.DB.Exec(command)
	if err != nil {
		log.Fatal("Unable to create table", err)
	}

	fmt.Println("set up database successfully")
}
