package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Error struct {
	Message string `json:"message"`
}

func getHandler(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func createExpenseHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount,note,tags) values ($1, $2, $3, $4)  RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))

	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, exp)

}

func updateHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "update")
}

func getAllHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "getAll")
}

var db *sql.DB

func main() {
	// os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", "postgres://rudwxbbq:UX7e7F375OJMZDzvtMd5BWlwenaRM0mv@tiny.db.elephantsql.com/rudwxbbq")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	command := "CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);"

	_, err = db.Exec(command)
	if err != nil {
		log.Fatal("Unable to create table", err)
	}

	e := echo.New()

	e.GET("/expenses", getAllHandler)
	e.GET("/expenses/:id", getHandler)
	e.POST("/expenses", createExpenseHandler)
	e.PUT("/expenses/:id", updateHandler)

	log.Fatal(e.Start(":2565"))
}
