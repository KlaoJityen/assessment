package main

import (
	"log"
	"net/http"

	"github.com/KlaoJityen/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func VerifyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != "November 10, 2009" {
			return c.JSON(http.StatusUnauthorized, "Unauthorized Access")
		}
		return next(c)
	}
}

func main() {

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(VerifyAuth)

	e.GET("/expenses", expense.GetExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.POST("/expenses", expense.CreateExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)

	log.Fatal(e.Start(":2565"))
}
