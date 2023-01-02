package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenseHandler(c echo.Context) error {

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
