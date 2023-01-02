package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpensesHandler(c echo.Context) error {
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

func GetExpenseHandler(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "ID is invalid"})
	}

	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't prepare query expense statement:" + err.Error()})
	}

	row := stmt.QueryRow(id)

	exp := Expense{}
	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Error{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't scan expenses:" + err.Error()})
	}
}
