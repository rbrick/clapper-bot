package main

import (
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	clapEmoji rune = '👏'
	openEmoji rune = '👐'
)

var (
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			// free advertisement
			tgbotapi.NewInlineKeyboardButtonURL("ClapperBot", "https://github.com/rbrick/clapper-bot"),
		),
	)
)

type QueuedClap struct {
	Timestamp  int64
	OldMessage *string
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true
	log.Printf("[ClapperBot] ClapperBot started...")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan, err := bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Fatal(err)
	}

	// message id -> timestamp & old message
	animatedMessages := map[string]*QueuedClap{}

	go func() {

		for {
			for k, v := range animatedMessages {
				expired := false
				if time.Now().Unix()-v.Timestamp >= 30 {
					delete(animatedMessages, k)
					expired = true
				}

				old := *v.OldMessage
				new := strings.Map(func(r rune) rune {
					// we want to end with a clap.
					if (r == openEmoji || r == clapEmoji) && expired {
						return clapEmoji
					} else if r == openEmoji {
						return clapEmoji
					} else if r == clapEmoji {
						return openEmoji
					}
					return r
				}, old)

				editMessage := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						InlineMessageID: k,
						ReplyMarkup:     &keyboard,
					},
					Text: new,
				}

				bot.Send(editMessage)

				if !expired {
					animatedMessages[k] = &QueuedClap{
						Timestamp:  v.Timestamp,
						OldMessage: &new,
					}
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()

	for update := range updateChan {
		if update.ChosenInlineResult != nil {
			result := update.ChosenInlineResult
			animatedMessages[result.InlineMessageID] = queueClappedString(result.Query)
		} else if update.InlineQuery != nil {
			query := update.InlineQuery
			clapifiedMessage := clapify(query.Query)
			log.Printf("[ClapperBot] Message: [%s], User: [%s]\n", query.Query, query.From.UserName)
			log.Printf("[ClapperBot] Clapified message [%s]\n", clapifiedMessage)

			article := tgbotapi.NewInlineQueryResultArticle(query.ID, "Send", clapifiedMessage)

			article.ReplyMarkup = &keyboard

			article.Description = "Clapified: " + clapifiedMessage
			results := []interface{}{article}

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: query.ID,
				IsPersonal:    false,
				Results:       results,
			}

			bot.AnswerInlineQuery(inlineConfig)
		}
	}
}

func queueClappedString(s string) *QueuedClap {
	s = clapify(s)
	return &QueuedClap{
		Timestamp:  time.Now().Unix(),
		OldMessage: &s,
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
