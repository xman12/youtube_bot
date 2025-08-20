package youtube

import (
	"database/sql"
	"fmt"
	"time"
)

func AddRequest(db *sql.DB, userID int, msgFromUser string, youtubeCode string, isShort bool, messageID int)  {
	sqlStatement := `INSERT INTO request_logs (user_id, message, video_code, is_short, created_at, message_id) 
			VALUES (?, ?, ?, ?, ?, ?)`

	t := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(sqlStatement, userID, msgFromUser, youtubeCode, isShort, t, messageID)
	if err != nil {
		fmt.Println(err.Error())
	}
}
