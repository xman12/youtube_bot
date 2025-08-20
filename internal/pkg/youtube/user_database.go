package youtube

import (
	"database/sql"
	"fmt"
	"time"
)

func GetUser(db *sql.DB, chatID int64, userID int) int  {
	rows, err := db.Query("select id from users where chat_id=? limit 1", chatID)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// обработка запросов с бд
	for rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return userID
}

func GetActiveUser(db *sql.DB, userID int, active int) int  {
	rows, err := db.Query("select active from users where id=? limit 1", userID)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// обработка запросов с бд
	for rows.Next() {
		err := rows.Scan(&active)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return active
}

func CreateUser(db *sql.DB, chatID int64, firstName string, userName string) int  {
	var userID int
	sqlStatementUser := `INSERT INTO users (chat_id, name, login, created_at) VALUES (?, ?, ?, ?)`
	t := time.Now().Format("2006-01-02 15:04:05")

	err := db.QueryRow(sqlStatementUser, chatID, firstName, userName, t)
	if err != nil {
		//return
	}

	rows, errQuery := db.Query("select id from users where chat_id=? limit 1", chatID)
	if errQuery != nil {
		fmt.Println(errQuery.Error())
	}
	defer func(rows *sql.Rows) {
		errClose := rows.Close()
		if errClose != nil {

		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return userID
}
