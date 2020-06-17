package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Configuration struct {
	BotApi                string `json:"bot_token"`
	MongoConnectionString string
}

func main() {
	// Carregar as variáveis do arquivo config.json
	var configuration Configuration
	byteFile, _ := ioutil.ReadFile("./config.dev.json")
	json.Unmarshal(byteFile, &configuration)
	configuration.MongoConnectionString = os.Getenv("MONGO_URL")

	fmt.Println(configuration)

	// Inicialização do bot
	bot, err := tgbotapi.NewBotAPI(configuration.BotApi)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}