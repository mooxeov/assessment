package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func TestCreateExpenseHandler(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var exp Expense

	res := request(http.MethodPost, uri("expense"), body)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, exp.ID)
	assert.Equal(t, "strawberry smoothie", exp.Title)
	assert.Equal(t, 79, exp.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
	assert.Equal(t, []string([]string{"food", "beverage"}), (exp.Tags))
}

func TestSelectExpenseWithIDHandler(t *testing.T) {
	var exp Expense

	res := request(http.MethodGet, uri("expense/2"), nil)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, exp.ID)
	assert.Equal(t, "strawberry smoothie", exp.Title)
	assert.Equal(t, 79, exp.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
	assert.Equal(t, []string([]string{"food", "beverage"}), (exp.Tags))
}

func TestUpdateExpenseHandler(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title":"apple smoothie",
		"amount":89,
		"note":"no discount", 
		"tags":["beverage"]
	}`)
	var exp Expense

	res := request(http.MethodPut, uri("expense/1"), body)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, exp.ID)
	assert.Equal(t, "apple smoothie", exp.Title)
	assert.Equal(t, 89, exp.Amount)
	assert.Equal(t, "no discount", exp.Note)
	assert.Equal(t, []string([]string{"beverage"}), (exp.Tags))
}

func TestSelectExpenseHandler(t *testing.T) {
	var exp [2]Expense

	res := request(http.MethodGet, uri("expense"), nil)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, exp[0].ID)
	assert.Equal(t, "apple smoothie", exp[0].Title)
	assert.Equal(t, 89, exp[0].Amount)
	assert.Equal(t, "no discount", exp[0].Note)
	assert.Equal(t, []string([]string{"beverage"}), (exp[0].Tags))

	assert.NotEqual(t, 0, exp[1].ID)
	assert.Equal(t, "iPhone 14 Pro Max 1TB", exp[1].Title)
	assert.Equal(t, 66900, exp[1].Amount)
	assert.Equal(t, "birthday gift from my love", exp[1].Note)
	assert.Equal(t, []string([]string{"gadget"}), (exp[1].Tags))
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
