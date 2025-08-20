package youtube

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"net/http"
	"os"
	"regexp"
)

func IsUrl(url string) bool {

	matched, _ := regexp.MatchString(`(.*)youtu(.*)`, url)
	if matched {
		return true
	}

	return false
}

func GetCode(url string) string {

	re, _ := regexp.Compile(`(be\/|v=)(.{11})`)
	res := re.FindString(url)

	replaceRe := regexp.MustCompile(`be\/|v=`)
	replaceResult := replaceRe.ReplaceAllString(res, "")

	return replaceResult
}

func GetShortCode(url string) string {
	re, _ := regexp.Compile(`shorts\/(.{11})`)
	res := re.FindString(url)

	replaceRe := regexp.MustCompile(`shorts\/`)
	replaceResult := replaceRe.ReplaceAllString(res, "")

	return replaceResult
}

func NumberOfPercent(percent string) string {

	re, _ := regexp.Compile(`([0-9]*\.[0-9]+)\%`)
	res := re.FindString(percent)

	replaceRe := regexp.MustCompile(`%`)
	replaceResult := replaceRe.ReplaceAllString(res, "")

	return replaceResult
}

func KeyboardForMedia(chatID int64, youtubeCode string) tgbotapi.InlineKeyboardMarkup {
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скачать видео", fmt.Sprintf("video|%d|%s", chatID, youtubeCode)),
			tgbotapi.NewInlineKeyboardButtonData("Скачать превью", fmt.Sprintf("preview|%d|%s", chatID, youtubeCode)),
		),
	)

	return numericKeyboard
}

// DownloadPhoto скачка фото превью/**
func DownloadPhoto(path string, youtubeCode string) {
	out, _ := os.Create(fmt.Sprintf("%s/%s.jpg", path, youtubeCode))
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)
	resp, _ := http.Get(fmt.Sprintf("https://i.ytimg.com/vi/%s/maxresdefault.jpg", youtubeCode))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	_, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetInfoForMessage() string {
	return "С уважением проект @GetYoutubeMp3Bot. Наш канал с музыкой: https://t.me/youtubemp3botchannel\n"
}
