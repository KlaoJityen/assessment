package main

import (
	"log"

	"github.com/KlaoJityen/assessment/expense"
	"github.com/labstack/echo/v4"
)

func main() {

	expense.InitDB()

	e := echo.New()

	e.GET("/expenses", expense.GetExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.POST("/expenses", expense.CreateExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)

	log.Fatal(e.Start(":2565"))
}
