package app

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	_ "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"youtube_bot/internal/pkg/youtube"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	counterNewUsers = promauto.NewCounter(prometheus.CounterOpts{
		Name: "register_users",
		Help: "Counting the total number of registrations handled",
	})

	counterRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "requests",
		Help: "Counting the total number of requests handled",
	})

	counterVideoLoads = promauto.NewCounter(prometheus.CounterOpts{
		Name: "video_loads",
		Help: "Counting the total number of video loads handled",
	})

	countVideoLoadsError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "video_loads_error",
		Help: "Counting the total number of video loads error",
	})
)

func Run() {

	srv := http.NewServeMux()
	srv.Handle("/metrics", promhttp.Handler())

	go func() {

		// Ниже запуск бота
		// Используем локальный Telegram API сервер если указан TELEGRAM_API_URL
		var bot *tgbotapi.BotAPI
		if apiURL := os.Getenv("TELEGRAM_API_URL"); apiURL != "" {
			bot, _ = tgbotapi.NewBotAPIWithAPIEndpoint(os.Getenv("TELEGRAM_API_TOKEN"), apiURL+"/bot%s/%s")
		} else {
			bot, _ = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
		}
		bot.Debug = true

		dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_LOGIN"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"))

		db, err := sql.Open("mysql", dbSource)

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {

			}
		}(db)

		if err != nil {
			panic(err)
		}

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30
		updates := bot.GetUpdatesChan(updateConfig)

		proxyKey := youtube.ProxyOption{
			Key:  0,
			Path: os.Getenv("PATH_TO_PROXY"),
		}

		// Let's go through each update that we're getting from Telegram.
		for update := range updates {
			var chatID int64
			if update.CallbackQuery == nil {
				chatID = update.Message.Chat.ID
				if chatID < 0 {
					continue
				}
			}

			if update.CallbackQuery != nil {
				data := update.CallbackQuery.Data
				explodeData := strings.Split(data, "|")
				chatID, _ := strconv.ParseInt(explodeData[1], 10, 64)
				youtubeCode := explodeData[2]
				if explodeData[0] == "video" {

					var fileID string
					fileID, _ = youtube.GetVideoVideoFile(db, youtubeCode)

					if fileID != "" {
						//msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "Видео скачивается, нужно немного подождать")

						msg := tgbotapi.NewVideo(chatID, tgbotapi.FileID(fileID))
						msg.Caption = youtube.GetInfoForMessage()
						_, _ = bot.Send(msg)
						continue
					}

					msg := tgbotapi.NewMessage(chatID, "Идет загрузка видео")
					result, _ := bot.Send(msg)

					go func() {
						videoFile := fmt.Sprintf("%s/%s.mp4", os.Getenv("PATH_TO_LOAD_VIDEO"), youtubeCode)
						videoFileWebm := fmt.Sprintf("%s/%s.webm", os.Getenv("PATH_TO_LOAD_VIDEO"), youtubeCode)
						videoFileMkv := fmt.Sprintf("%s/%s.mkv", os.Getenv("PATH_TO_LOAD_VIDEO"), youtubeCode)

						if false == Exists(videoFile) && false == Exists(videoFileWebm) && false == Exists(videoFileMkv) {
							msgEdit := tgbotapi.NewEditMessageText(chatID, result.MessageID, "Идет загрузка...")

							youtube.DownloadVideo(
								youtubeCode,
								os.Getenv("PATH_TO_LOAD_VIDEO"),
								os.Getenv("PATH_TO_LOAD_AUDIO"),
								false,
								*bot,
								msgEdit,
								&proxyKey,
								result.MessageID,
								db,
								countVideoLoadsError,
							)
						}

						if true == Exists(videoFileMkv) {
							GetMp4(videoFileMkv)
							videoFile = fmt.Sprintf("%s/%s.mkv.mp4", os.Getenv("PATH_TO_LOAD_VIDEO"), youtubeCode)
						}

						if false == Exists(videoFile) && true == Exists(videoFileWebm) {
							args := fmt.Sprintf("ffmpeg -fflags +genpts -i %s -r 24 %s", videoFileWebm, videoFile)
							cmd := exec.Command("bash", "-c", args)
							stderr, _ := cmd.StdoutPipe()
							err := cmd.Start()
							if err != nil {
								return
							}
							scanner := bufio.NewScanner(stderr)
							scanner.Split(bufio.ScanWords)
							for scanner.Scan() {
								m := scanner.Text()
								fmt.Println(m)
							}
						}

						msgMedia := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(videoFile))
						msgMedia.Caption = youtube.GetInfoForMessage()

						result, err := bot.Send(msgMedia)
						if err != nil {
							fmt.Println(err.Error())
						}
						youtube.UpdateVideoField(db, youtubeCode, result.Video.FileID)

					}()

				}

				if explodeData[0] == "preview" {
					_, _ = directoryExistsAndCreate(os.Getenv("PATH_TO_LOAD_IMG"))
					youtube.DownloadPhoto(os.Getenv("PATH_TO_LOAD_IMG"), explodeData[2])
					imgPath := fmt.Sprintf("%s/%s.jpg", os.Getenv("PATH_TO_LOAD_IMG"), explodeData[2])
					file := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imgPath))
					file.Caption = youtube.GetInfoForMessage()
					_, err := bot.Send(file)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}

			if update.Message == nil {
				continue
			}

			// Count requests
			counterRequests.Inc()

			var userID int
			userID = 0
			active := 0
			audioFileID := ""

			// получение пользователя
			userID = youtube.GetUser(db, update.Message.Chat.ID, userID)

			if userID == 0 {
				// создание пользователя
				userID = youtube.CreateUser(db, update.Message.Chat.ID, update.Message.Chat.FirstName, update.Message.Chat.UserName)
				// Counter new users
				counterNewUsers.Inc()
			}

			active = youtube.GetActiveUser(db, userID, active)
			if 0 == active {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Невозможно обработать запроc")
				_, _ = bot.Send(msg)
				continue
			}

			msgForUser := update.Message.Text
			msgFromUser := update.Message.Text
			isShort := false
			youtubeCode := youtube.GetCode(msgForUser)
			var length = len([]rune(youtubeCode))
			if length == 0 {
				youtubeCode = youtube.GetShortCode(msgForUser)
				length = len([]rune(youtubeCode))
				isShort = true
			}

			if youtube.IsUrl(update.Message.Text) {
				msgForUser = fmt.Sprintf("Код видео: %s, скоро мы начнем скачивание", youtubeCode)
			} else {
				msgForUser = "Приветствую! Данный бот предназначен для получения mp3 файлов из роликов на Youtube.com\nШли ссылку а я ее обработаю 😏"
			}

			// поиск уже имеющегося ролика, если он находится то сразу отдаем файл
			audioFileID = youtube.GetVideoAudioFile(db, youtubeCode)

			if audioFileID != "" {

				var fileID string
				fileID = audioFileID

				numericKeyboard := youtube.KeyboardForMedia(chatID, youtubeCode) //кнопки ресурс
				msgMedia := tgbotapi.NewAudio(chatID, tgbotapi.FileID(fileID))
				msgMedia.Caption = youtube.GetInfoForMessage()
				msgMedia.ReplyMarkup = numericKeyboard
				_, err := bot.Send(msgMedia)
				if err != nil {
					fmt.Println(err.Error())
				}
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgForUser)
			result, err2 := bot.Send(msg)
			if err2 != nil {
				//panic(err)
			}

			// запуск обработки и скачки роликов
			update := update
			go func() {
				if 0 < length {
					// считаем какое количество запросов идет на скачку видео
					counterVideoLoads.Inc()
					msgEdit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, result.MessageID, "Идет загрузка")

					youtube.AddRequest(db, userID, msgFromUser, youtubeCode, isShort, result.MessageID)

					// проверка наличия директорий, и если их нет, то создаем
					_, _ = directoryExistsAndCreate(os.Getenv("PATH_TO_LOAD_AUDIO"))
					_, _ = directoryExistsAndCreate(os.Getenv("PATH_TO_LOAD_VIDEO"))

					// добавляем и скачиваем видео
					if audioFileID == "" {
						youtube.AddVideo(db, youtubeCode)
						youtube.DownloadVideo(
							youtubeCode,
							os.Getenv("PATH_TO_LOAD_VIDEO"),
							os.Getenv("PATH_TO_LOAD_AUDIO"),
							isShort,
							*bot,
							msgEdit,
							&proxyKey,
							result.MessageID,
							db,
							countVideoLoadsError,
						)
					}
				}
			}()
		}
	}()

	if err := http.ListenAndServe(":8090", srv); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetMp4(videoFile string) {
	args := fmt.Sprintf("ffmpeg -i %s -crf 18 %s.mp4", videoFile, videoFile)
	cmd := exec.Command("bash", "-c", args)

	stderr, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func directoryExistsAndCreate(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // Папка существует
	}
	if os.IsNotExist(err) {
		err := CreateDirectory(path)
		if err != nil {
			return false, err
		}
		return false, nil // Папка не существует
	}
	return false, err // Другая ошибка
}

func CreateDirectory(path string) error {
	err := os.Mkdir(path, 0777) // 0755 - права доступа (chmod)
	if err != nil {
		return err
	}
	return nil
}
