package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func getHandler(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func createHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "create")
}

func updateHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "update")
}

func getAllHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "getAll")
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	e := echo.New()

	e.GET("/users", getAllHandler)
	e.GET("/users/:id", getHandler)
	e.POST("/users", createHandler)
	e.PUT("/users/:id", updateHandler)

	log.Fatal(e.Start(":2565"))
}
