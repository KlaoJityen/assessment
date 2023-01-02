package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
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

func main() {
	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))

	e := echo.New()

	e.GET("/users", getAllHandler)
	e.GET("/users/:id", getHandler)
	e.POST("/users", createHandler)
	e.PUT("/users/:id", updateHandler)

	log.Fatal(e.Start(":2565"))
}
