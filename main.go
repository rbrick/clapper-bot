package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const clapEmoji = '👏'

func main() {
	bot, err := tgbotapi.NewBotAPI("395173942:AAHpz_VEacfdIdEDhGMDEHl5oZLcr6q232A")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("ClapperBot started...")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan, err := bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Fatal(err)
	}

	for update := range updateChan {
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
