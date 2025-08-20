package youtube

import (
	"bufio"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

type Process struct {
	processId int64
	name      string
}

func DownloadVideo(
	url string,
	pathLoad string,
	pathAudio string,
	isShort bool,
	bot tgbotapi.BotAPI,
	msg tgbotapi.EditMessageTextConfig,
	proxyOption *ProxyOption,
	messageId int,
	db *sql.DB,
	countVideoLoadsError prometheus.Counter,
) {

	urlDownload := url
	if isShort {
		urlDownload = fmt.Sprintf("https://www.youtube.com/shorts/%s", url)
	}

	proxy := GetProxy(proxyOption)
	fmt.Println(proxy.GetProxyString())
	args := fmt.Sprintf("yt-dlp --no-mtime --proxy %s -o \"%s/%%(id)s.%%(ext)s\" \"%s\" --merge-output-format=\"mp4/mkv\" -f w", proxy.GetProxyString(), pathLoad, urlDownload)
	//fmt.Println(args)
	cmd := exec.Command("bash", "-c", args)

	stderr, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return
	}
	step := 1
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	canDownload := false
	numberPercentFull := 0.00
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		matched, _ := regexp.MatchString(`ERROR`, m)
		if false == matched {
			canDownload = true
			numberPercent := NumberOfPercent(m)

			percentInt, err := strconv.ParseFloat(numberPercent, 64)
			if err != nil {

			}
			if numberPercentFull < percentInt {
				numberPercentFull = percentInt
			}

			fmt.Println(numberPercent)
			fmt.Println(percentInt)
			step++

			// каждые 5 делений показываем статус в боте о этапах скачки файла, чем меньше шаг тем чаще шлет обновления
			if 100 == step || "100" == numberPercent {
				step = 1

				msg.Text = fmt.Sprintf("Идет загрука: %.2f %%", numberPercentFull)

				if _, err := bot.Send(msg); err != nil {
					//panic(err)
				}

				if "100" == numberPercent {

					msg.Text = "Обрабатываем данные на сервере..."
					_, _ = bot.Send(msg)
				}

				fmt.Println(m)
			}
		} else {
			canDownload = false
		}
	}

	if canDownload == false {
		countVideoLoadsError.Inc()
		msg.Text = "К сожалению я не могу загрузить данный ролик"
		if _, err := bot.Send(msg); err != nil {
			//panic(err)
		}
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	msg.Text = "Выгружаем аудио дорожку из видеофайла"
	_, _ = bot.Send(msg)

	videoFile := fmt.Sprintf("%s/%s.mp4", pathLoad, url)
	audioFile := fmt.Sprintf("%s/%s.mp3", pathAudio, url)

	// проверка существования mp4 если нет то это webm
	if false == Exists(videoFile) {
		videoFile = fmt.Sprintf("%s/%s.webm", pathLoad, url)

		if false == Exists(videoFile) {
			videoFile = fmt.Sprintf("%s/%s.mkv", pathLoad, url)
			GetMp4(videoFile)
		}
	}
	GetMp3(videoFile, audioFile)

	msg.Text = "Выгружаем аудио дорожку из видеофайла"
	_, _ = bot.Send(msg)

	msgAudio := tgbotapi.NewAudio(msg.ChatID, tgbotapi.FilePath(audioFile))
	msgAudio.ReplyMarkup = KeyboardForMedia(msg.ChatID, url) //кнопки ресурс
	msgAudio.Caption = GetInfoForMessage()
	result, err := bot.Send(msgAudio)
	if err != nil {
		return
	}

	// обновляем в БД данные
	UpdateAudioField(db, url, result.Audio.FileID) // обновляем в БД данные
	return
}

func GetMp3(
	videoFile string,
	audioFile string,
) {

	args := fmt.Sprintf("ffmpeg -i %s %s", videoFile, audioFile)
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

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
