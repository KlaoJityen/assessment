package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(c echo.Context) error {
	fmt.Println("inside handler")
	return c.JSON(http.StatusOK, "hello world")
}

func main() {
	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))

	e := echo.New()

	e.GET("/users", Handler)

	log.Fatal(e.Start(":2565"))
}
