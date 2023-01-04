package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/KlaoJityen/assessment/expense"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func TestGetExpenses(t *testing.T) {
	reqBody := `{"id":1,"title":"groceries","amount":50,"note":"","tags":["food", "necessities"]}`

	req, err := http.NewRequest(http.MethodPost, uri("expenses"), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")

	var resp *http.Response
	client := http.Client{}
	_, err = client.Do(req)
	assert.NoError(t, err)

	req, err = http.NewRequest(http.MethodGet, uri("expenses"), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")

	client = http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBuffer []byte
	byteBuffer, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var exp []expense.Expense
	err = json.Unmarshal(byteBuffer, &exp)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Greater(t, len(exp), 0)
}

func TestGetExpenseByID(t *testing.T) {
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, uri("expenses/1"), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")

	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBuffer []byte
	byteBuffer, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var exp expense.Expense
	err = json.Unmarshal(byteBuffer, &exp)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, exp.ID, 1)
	assert.Equal(t, "apple smoothie", exp.Title)
	assert.Equal(t, float64(10), exp.Amount)
}

func TestCreateExpense(t *testing.T) {
	reqBody := `{"title":"java smoothie","amount":50,"note":"write one, bug anywhere","tags":["food", "necessities"]}`

	req, err := http.NewRequest(http.MethodPost, uri("expenses"), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")

	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBuffer []byte
	byteBuffer, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var exp expense.Expense
	err = json.Unmarshal(byteBuffer, &exp)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEqual(t, 0, exp.ID)
	assert.Equal(t, "java smoothie", exp.Title)
	assert.Equal(t, float64(50), exp.Amount)
	assert.Equal(t, "write one, bug anywhere", exp.Note)
	assert.Equal(t, []string{"food", "necessities"}, exp.Tags)

}

func TestUpdateExpense(t *testing.T) {
	reqBody := `{
        "id": 1,
        "title": "apple smoothie",
        "amount": 10,
        "note": "no discount",
        "tags": [
            "beverage"
        ]
    }`

	req, err := http.NewRequest(http.MethodPut, uri("expenses/1"), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", "November 10, 2009")

	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBody []byte
	byteBody, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var exp expense.Expense
	err = json.Unmarshal(byteBody, &exp)
	assert.NoError(t, err)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, exp.ID)
	assert.Equal(t, "apple smoothie", exp.Title)
	assert.Equal(t, float64(10), exp.Amount)
	assert.Equal(t, "no discount", exp.Note)
	assert.Equal(t, []string{"beverage"}, exp.Tags)

}
