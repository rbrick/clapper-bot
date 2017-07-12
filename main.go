package clapperbot

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://clapper-bot.appspot.com:8443/"+bot.Token, "cert.pem"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[ClapperBot] ClapperBot started...")

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		if update.InlineQuery != nil {
			query := update.InlineQuery
			clapifiedMessage := clapify(query.Query)

			log.Printf("[ClapperBot] Message: [%s], User: [%s]\n", query.Query, query.From.UserName)
			log.Printf("[ClapperBot] Clapified message [%s]\n", clapifiedMessage)

			article := tgbotapi.NewInlineQueryResultArticle(query.ID, "Send", clapifiedMessage)
			article.Description = "Clapified: " + clapifiedMessage
			results := []interface{}{article}

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: query.ID,
				Results:       results,
			}

			bot.AnswerInlineQuery(inlineConfig)
		} else {
			log.Println("[ClapperBot] Inline Mode not used. :////")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Must 👏 use 👏 inline."))
		}
	}
}

func clapify(s string) string {
	var clapEmoji = '👏'

	strs := strings.Split(s, " ")
	ns := ""

	for i, v := range strs {
		ns += v
		if i == len(strs)-1 {
			continue
		}
		ns += " " + string(clapEmoji) + " "
	}

	return ns
}
