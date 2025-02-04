package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateExpenseHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	row := handler.DB.QueryRow("INSERT INTO expenses (title, amount,note,tags) values ($1, $2, $3, $4)  RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))

	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, exp)

}
