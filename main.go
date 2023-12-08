package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Default()

	logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	logger.Handler()

	config := slogecho.Config{
		WithRequestBody:    true,
		WithResponseBody:   true,
		WithRequestHeader:  true,
		WithResponseHeader: true,
	}

	h := Handler{bot: bot}
	e := echo.New()
	e.Use(slogecho.NewWithConfig(logger, config))
	e.Use(middleware.Recover())
	e.POST("/", h.HandleUpdate)
	err = e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		logger.Error(err.Error())
	}
}

type Handler struct {
	bot *tgbotapi.BotAPI
}

func (h *Handler) HandleUpdate(c echo.Context) error {
	logger.Info(fmt.Sprintf("c.Request().RequestURI: %v\n", c.Request().RequestURI))
	logger.Info(fmt.Sprintf("c.Request().Method: %v\n", c.Request().Method))
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		log.Print("Cannot bind update", err)
		return err
	}
	HandleMessage(update, h.bot)
	return c.JSON(200, nil)
}

func HandleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := update.Message
	if message == nil {
		logger.Info("Message not founded")
		return
	}
	userName := update.Message.From.UserName
	// ID чата/диалога.
	// Может быть идентификатором как чата с пользователем
	// (тогда он равен UserID) так и публичного чата/канала
	ChatID := update.Message.Chat.ID
	// Текст сообщения
	Text := update.Message.Text

	logger.Info("Handle message", slog.String("userName", userName), slog.String("chat", update.Message.Chat.Title))

	if res, ok := convertInstReel(Text); ok {
		messageId := update.Message.MessageID
		delMsg := tgbotapi.NewDeleteMessage(ChatID, messageId)
		newMsg := tgbotapi.NewMessage(ChatID, fmt.Sprintf("Отправил вот этот человек: @%v \n %v", userName, res))
		bot.Send(delMsg)
		bot.Send(newMsg)
	}
}
