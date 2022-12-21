package expense

import (
	"database/sql"
	"log"
)

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
