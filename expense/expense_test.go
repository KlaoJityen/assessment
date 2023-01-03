package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetExpense(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"ID", "Title", "Amount", "Note", "Tags"}).
		AddRow("1", "apple smoothie", 89, "no discount", pq.Array(&[]string{"beverage"}))

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	handler = &Handler{db}

	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").ExpectQuery().
		WithArgs(1).WillReturnRows(newsMockRows)

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := `{"id":1,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`
	// Act
	err = GetExpenseHandler(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	newsMockRows := sqlmock.NewRows([]string{"ID", "Title", "Amount", "Note", "Tags"}).
		AddRow("1", "apple smoothie", 89, "no discount", pq.Array(&[]string{"beverage"})).
		AddRow("2", "coconut smoothie", 89, "no discount", pq.Array(&[]string{"beverage"}))

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC").ExpectQuery().WillReturnRows(newsMockRows)

	handler = &Handler{db}

	c := e.NewContext(req, rec)
	expected := `[{"id":1,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]},{"id":2,"title":"coconut smoothie","amount":89,"note":"no discount","tags":["beverage"]}]`

	// Act
	err = GetExpensesHandler(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestCreateExpense(t *testing.T) {
	// Arrange

	mockBody := strings.NewReader(`{"title":"groceries","amount":50,"note":"","tags":["food", "necessities"]}`)
	exp := Expense{ID: 1, Title: "groceries", Amount: 50, Note: "", Tags: []string{"food", "necessities"}}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", mockBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"id"}).AddRow(exp.ID)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	mock.ExpectQuery("INSERT INTO expenses").WithArgs(exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)).
		WillReturnRows(newsMockRows)

	handler = &Handler{db}
	c := e.NewContext(req, rec)

	// Act
	err = CreateExpenseHandler(c)

	// Assertions
	expected := `{"id":1,"title":"groceries","amount":50,"note":"","tags":["food","necessities"]}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestUpdateExpense(t *testing.T) {
	// Arrange

	mockBody := strings.NewReader(`{"id":1,"title":"groceries","amount":50,"note":"","tags":["food", "necessities"]}`)
	exp := Expense{ID: 1, Title: "groceries", Amount: 50, Note: "", Tags: []string{"food", "necessities"}}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/1", mockBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewResult(1, 1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	mock.ExpectPrepare("UPDATE expenses").ExpectExec().
		WithArgs(exp.ID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)).WillReturnResult(newsMockRows)

	handler = &Handler{db}
	c := e.NewContext(req, rec)

	// Act
	err = UpdateExpenseHandler(c)

	// Assertions
	expected := `{"id":1,"title":"groceries","amount":50,"note":"","tags":["food","necessities"]}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
