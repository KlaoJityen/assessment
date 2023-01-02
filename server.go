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

func getExpensesHandler(c echo.Context) error {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't prepare query for expenses statement:" + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't query expenses:" + err.Error()})
	}

	expenses := []Expense{}

	for rows.Next() {
		expense := Expense{}
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Error{Message: "can't scan row:" + err.Error()})
		}
		expenses = append(expenses, expense)
	}

	return c.JSON(http.StatusOK, expenses)
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

func updateExpenseHandler(c echo.Context) error {

	exp := Expense{}
	err := c.Bind(&exp)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount , note, tags")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't prepare query expense statment:" + err.Error()})
	}

	if _, err := stmt.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't scan expenses:" + err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}

func getExpenseHandler(c echo.Context) error {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't prepare query expense statement:" + err.Error()})
	}

	row := stmt.QueryRow(id)

	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Error{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't scan expenses:" + err.Error()})
	}
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

	e.GET("/expenses", getExpensesHandler)
	e.GET("/expenses/:id", getExpenseHandler)
	e.POST("/expenses", createExpenseHandler)
	e.PUT("/expenses/:id", updateExpenseHandler)

	log.Fatal(e.Start(":2565"))
}
