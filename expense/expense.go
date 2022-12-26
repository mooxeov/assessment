package expense

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func CreateExpenseHandler(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values($1,$2,$3,$4) RETURNING id,title,amount,note,tags", e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (*pq.StringArray)(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func GetExpenseWithIDHandler(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	var e Expense

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("SELECT id,title, amount, note, tags FROM expenses WHERE id = $1", rowID)

	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (*pq.StringArray)(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func UpdateExpenseHandler(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	var e Expense
	err = c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1 RETURNING id,title,amount,note,tags", rowID, e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (*pq.StringArray)(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

var db *sql.DB

func InitDB() {
	url := "postgres://hbaivmlu:dDt9ul_KZn00fZuzx7Wol8Qza2cmvob1@satao.db.elephantsql.com/hbaivmlu"
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT, amount INT, note TEXT, tags TEXT [])
	`

	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

}
