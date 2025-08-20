package youtube

import (
	"database/sql"
	"fmt"
)

func GetVideoAudioFile(db *sql.DB, videoCode string) string {

	var audioFileID string

	videoRow, errVideo := db.Query("SELECT audio_file_id from videos where code=?", videoCode)
	if errVideo != nil {

	}

	defer func(videoRow *sql.Rows) {
		err := videoRow.Close()
		if err != nil {

		}
	}(videoRow)

	for videoRow.Next() {
		err := videoRow.Scan(&audioFileID)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return audioFileID
}


func GetVideoVideoFile(db *sql.DB, videoCode string) (string, error){

	var videoFileID sql.NullString

	videoRow, errVideo := db.Query("SELECT video_file_id from videos where code=?", videoCode)
	if errVideo != nil {

	}

	defer func(videoRow *sql.Rows) {
		err := videoRow.Close()
		if err != nil {

		}
	}(videoRow)

	for videoRow.Next() {
		err := videoRow.Scan(&videoFileID)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return videoFileID.String, nil
}

func UpdateAudioField(db *sql.DB, videoCode string, audioFileId string)  {
	sqlStatementVideo := `UPDATE videos SET audio_file_id=? where code=?`
	_, _ = db.Exec(sqlStatementVideo, audioFileId, videoCode)
}

func UpdateVideoField(db *sql.DB, videoCode string, videoFileId string)  {
	sqlStatementVideo := `UPDATE videos SET video_file_id=? where code=?`
	_, _ = db.Exec(sqlStatementVideo, videoFileId, videoCode)
}

func AddVideo(db *sql.DB, videoCode string)  {
	sqlStatementVideo := `INSERT INTO videos (code) VALUES (?)`
	_, _ = db.Exec(sqlStatementVideo, videoCode)
}
