package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: реализовать румтур
// TODO: поиск по ticker.description
// TODO: svg to png
// TODO: если #BABA?M5 - stockscores, #BABA! - finviz, #BABA?! - shortvalue
// TODO: если в сообщении пользователя только команда - удалять его после обработки
// TODO: README
// TODO: вернуть возврат ссылок для inline mode

func main() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
	)
	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}
	pref := tb.Settings{
		// URL:    "https://api.bots.mn/telegram/",
		Token:  token,
		Poller: webhook,
	}
	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	b.Handle(tb.OnQuery, func(q *tb.Query) {
		re := regexp.MustCompile("[^A-Za-z]")
		symbol := re.ReplaceAllString(q.Text, "")
		ticker := GetExactTicker(symbol)
		if ticker == nil {
			return
		}
		results := make(tb.Results, 1+len(ArticleCases)) // []tb.Result
		linkURL := fmt.Sprintf("https://ru.tradingview.com/symbols/%s", ticker.symbol)
		result := &tb.ArticleResult{
			Title:       ticker.symbol,
			Description: ticker.description,
			HideURL:     true,
			URL:         linkURL,
			ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticker.symbol), // from stockanalysis.com
		}
		result.SetContent(&tb.InputTextMessageContent{
			Text: fmt.Sprintf(`\#%s \- [%s](%s)`,
				ticker.symbol,
				escape(ticker.description),
				linkURL,
			),
			ParseMode:      tb.ModeMarkdownV2,
			DisablePreview: true,
		})
		result.SetResultID(ticker.symbol)
		results[0] = result
		for i, articleCase := range ArticleCases {
			linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
			title := articleCase.name
			if articleCase.hasGift {
				title += " 🎁"
			}
			result := &tb.ArticleResult{
				Title:       title,
				Description: ticker.symbol,
				HideURL:     true,
				URL:         linkURL,
			}
			result.SetContent(&tb.InputTextMessageContent{
				Text: fmt.Sprintf("/info %s %s",
					articleCase.name,
					ticker.symbol,
				),
				DisablePreview: true,
			})
			result.SetResultID(ticker.symbol + "=" + articleCase.name)
			results[i+1] = result
		}
		err = b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})
		if err != nil {
			log.Println(err)
		}
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		if strings.HasPrefix(m.Text, "/info ") {
			re := regexp.MustCompile(",")
			payload := re.ReplaceAllString(m.Payload, " ")
			arguments := strings.Split(payload, " ")
			symbols := arguments[1:]
			if len(symbols) == 0 {
				log.Println("Empty symbols")
				return
			}
			articleCaseName := arguments[0]
			articleCase := GetExactArticleCase(articleCaseName)
			if articleCase == nil {
				log.Println("Invalid command")
				return
			}
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
					symbol = symbol[1:]
				}
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					continue
				}
				if articleCaseName == "finviz.com" {
					linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
					screenshot := Screenshot(linkURL)
					photo := &tb.Photo{
						File: tb.FromReader(bytes.NewReader(screenshot)),
						Caption: fmt.Sprintf(
							`\#%s [%s](%s) by %s`,
							ticker.symbol,
							escape(articleCase.name),
							linkURL,
							getUserLink(m.Sender),
						),
					}
					_, err := b.Send(
						tb.ChatID(m.Chat.ID),
						photo,
						&tb.SendOptions{
							ParseMode: tb.ModeMarkdownV2,
						},
					)
					if err != nil {
						log.Println(err)
					}
				}
				if articleCase.imageURL != "" {
					imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol)
					linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
					photo := &tb.Photo{
						File: tb.FromURL(imageURL),
						Caption: fmt.Sprintf(
							`\#%s [%s](%s) by %s`,
							ticker.symbol,
							escape(articleCase.name),
							linkURL,
							getUserLink(m.Sender),
						),
					}
					_, err := b.Send(
						tb.ChatID(m.Chat.ID),
						photo,
						&tb.SendOptions{
							ParseMode: tb.ModeMarkdownV2,
						},
					)
					if err != nil {
						log.Println(err)
					}
				}
			}
			err := b.Delete(
				&tb.StoredMessage{
					MessageID: strconv.Itoa(m.ID),
					ChatID:    m.Chat.ID,
				},
			)
			if err != nil {
				log.Println(err)
			}
		}
	})
	b.Start()
}

func contains(slice []string, search string) bool {
	for _, element := range slice {
		if element == search {
			return true
		}
	}
	return false
}

func parseInt(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return result
}

func escape(s string) string {
	re := regexp.MustCompile("[.|-]")
	return re.ReplaceAllString(s, `\$0`)
}

// func deleteCommand(b *tb.Bot, m *tb.Message) {
// 	err := b.Delete(
// 		&tb.StoredMessage{
// 			MessageID: strconv.Itoa(m.ID),
// 			ChatID:    m.Chat.ID,
// 		},
// 	)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func sendInformer(b *tb.Bot, m *tb.Message, photo *tb.Photo) {
// 	_, err := b.Send(
// 		tb.ChatID(m.Chat.ID),
// 		photo,
// 		&tb.SendOptions{
// 			ParseMode: tb.ModeMarkdownV2,
// 		},
// 	)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func getUserLink(u *tb.User) string {
	if u.Username != "" {
		return fmt.Sprintf("@%s", u.Username)
	}
	return fmt.Sprintf("[%s](tg://user?id=%d)", u.FirstName, u.ID)
}
