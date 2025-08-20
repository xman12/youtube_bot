package youtube

import (
	"bufio"
	"context"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/telegram/uploader"
	"go.uber.org/zap"
	//"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/message/html"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh/terminal"
)

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	noSignUp

	phone string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadToChat(
	chatID int64,
	fileName string,
	pathLoad string,
	typeFile string,
	fileTitle string,
	messageId int,
	) {

	Run(func(ctx context.Context, log *zap.Logger) error {
		// Dispatcher handles incoming updates.
		dispatcher := tg.NewUpdateDispatcher()
		gaps := updates.New(updates.Config{
			Handler: dispatcher,
			Logger:  log.Named("gaps"),
		})
		opts := telegram.Options{
			Logger:        log,
			UpdateHandler: gaps,
			Middlewares: []telegram.Middleware{
				updhook.UpdateHook(gaps.Handle),
			},
		}

		client, _ := telegram.ClientFromEnvironment(opts)

		dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
			//log.Info("Message", zap.Any("message", update.Message))
			return nil
		})
		dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
			//log.Info("Message", zap.Any("message", update.Message))
			return nil
		})

		return client.Run(ctx, func(ctx context.Context) error {
			// Note: you need to be authenticated here.

			api := tg.NewClient(client)

			// Helper for uploading. Automatically uses big file upload when needed.
			u := uploader.NewUploader(api)
			//
			//// Helper for sending messages.
			sender := message.NewSender(api).WithUploader(u)
			filePath := fmt.Sprintf("%s/%s", pathLoad, fileName)
			// Uploading directly from path. Note that you can do it from
			// io.Reader or buffer, see From* methods of uploader.
			log.Info("Uploading file")
			upload, err := u.FromPath(ctx, filePath)
			if err != nil {
				return fmt.Errorf("upload %q: %w", filePath, err)
			}

			// Now we have uploaded file handle, sending it as styled message.
			// First, preparing message.
			document := message.UploadedDocument(upload,
				html.String(nil, fmt.Sprintf(`%d|%s|%s|%d`, chatID, typeFile, fileTitle, messageId)),
			)

			// You can set MIME type, send file as video or audio by using
			// document builder:
			if typeFile == "audio" {
				document.
					MIME("audio/mp3").
					Filename(fmt.Sprintf("%s.mp3", fileTitle)).
					Audio()
			} else {
				document.
					MIME("video/mp4").
					Filename(fmt.Sprintf("%s.mp4", fileTitle)).
					Video()
			}


			// Resolving target. Can be telephone number or @nickname of user,
			// group or channel.
			target := sender.ResolveDeeplink("https://t.me/youtube_ham")

			// Sending message with media.
			log.Info("Sending file")
			_, err = target.Media(ctx, document)
			if err != nil {
				//return fmt.Errorf("send: %w", err)
			}

			<-ctx.Done()
			return ctx.Err()
		})
	})
}

func Run(f func(ctx context.Context, log *zap.Logger) error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		//panic(err)
	}
	defer func() { _ = log.Sync() }()
	// No graceful shutdown.
	ctx := context.Background()
	err = f(ctx, log)
	//err != nil {
	//	log.Fatal("Run failed", zap.Error(err))
	//}
	// Done.
}
