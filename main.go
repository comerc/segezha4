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
		log.Println("****")
		log.Println(m.Text)
		log.Println(b.Me)
		log.Println(m.Via)
		log.Println(m.Chat)
		log.Println(m.Sender)
		log.Println("****")
		mode := ""
		if m.Via != nil && m.Via.ID == b.Me.ID {
			mode = "inline mode"
		} else if m.Chat.Type == tb.ChatPrivate {
			mode = "command mode"
		}
		if mode == "" {
			return
		}
		log.Println(mode)
		// а когда Send?
		if strings.HasPrefix(m.Text, "/info ") {
			params := strings.Split(m.Payload, " ")
			// TODO: много тикеров за один раз
			if len(params) < 2 {
				log.Println("Invalid params")
			}
			articleCaseName := params[0]
			tickerSymbol := params[1]
			articleCase := GetExactArticleCase(articleCaseName)
			if articleCaseName == "finviz.com" {
				linkURL := fmt.Sprintf(articleCase.linkURL, tickerSymbol)
				screenshot := Screenshot(linkURL)
				photo := &tb.Photo{
					File: tb.FromReader(bytes.NewReader(screenshot)),
					Caption: fmt.Sprintf(
						`\#%s [%s](%s)`,
						tickerSymbol,
						escape(articleCase.name),
						linkURL,
					),
				}
				b.Send(
					tb.ChatID(m.Chat.ID),
					photo,
					&tb.SendOptions{
						ParseMode: tb.ModeMarkdownV2,
					},
				)
			}
			if articleCase.imageURL != "" {
				imageURL := fmt.Sprintf(articleCase.imageURL, tickerSymbol)
				linkURL := fmt.Sprintf(articleCase.linkURL, tickerSymbol)
				photo := &tb.Photo{
					File: tb.FromURL(imageURL),
					Caption: fmt.Sprintf(
						`\#%s [%s](%s)`,
						tickerSymbol,
						escape(articleCase.name),
						linkURL,
					),
				}
				b.Send(
					tb.ChatID(m.Chat.ID),
					photo,
					&tb.SendOptions{
						ParseMode: tb.ModeMarkdownV2,
					},
				)
			}
		}
		// fmt.Println(m.Chat)
		// b.Send(m.Sender, )

		// err := b.Delete(
		// 	&tb.StoredMessage{
		// 		MessageID: strconv.Itoa(m.ID),
		// 		ChatID:    m.Chat.ID,
		// 	},
		// )
		// if err != nil {
		// 	log.Println(err)
		// }

		// b.Send(m.Sender, "hello world"+strconv.FormatInt(m.Chat.ID, 10))
	})
	// b.Handle(tb.OnChosenInlineResult, func(r *tb.ChosenInlineResult) {
	// 	// incoming inline queries
	// 	log.Println("====")
	// 	log.Println(r.MessageID)
	// 	log.Println(r.ResultID)
	// 	log.Println(r.Query)
	// 	log.Println(r.From.ID)
	// 	log.Println("====")
	// 	resultID := strings.Split(r.ResultID, "=")
	// 	if len(resultID) != 2 {
	// 		// TODO: empty message
	// 		return
	// 	}
	// 	tickerSymbol := resultID[0]
	// 	articleCaseName := resultID[1]
	// 	log.Println(articleCaseName)
	// 	log.Println(tickerSymbol)
	// 	// ticketName := r.ResultID
	// 	// TODO: to https://core.telegram.org/bots#deep-linking-example
	// 	to := tb.ChatID(parseInt(chatID))
	// 	// commands := make([]string, 0)
	// 	// for _, param := range strings.Split(r.Query, " ") {
	// 	// 	if strings.HasPrefix(param, "#") || strings.HasPrefix(param, "$") {
	// 	// 		continue
	// 	// 	}
	// 	// 	commands = append(commands, param)
	// 	// }
	// 	articleCase := GetExactArticleCase(articleCaseName)
	// 	if articleCaseName == "finviz.com" {
	// 		linkURL := fmt.Sprintf(articleCase.linkURL, tickerSymbol)
	// 		screenshot := Screenshot(linkURL)
	// 		photo := &tb.Photo{
	// 			File: tb.FromReader(bytes.NewReader(screenshot)),
	// 			Caption: fmt.Sprintf(
	// 				`\#%s [%s](%s)`,
	// 				tickerSymbol,
	// 				escape(articleCase.name),
	// 				linkURL,
	// 			),
	// 		}
	// 		b.Send(
	// 			to,
	// 			photo,
	// 			&tb.SendOptions{
	// 				ParseMode: tb.ModeMarkdownV2,
	// 			},
	// 		)
	// 	}
	// 	if articleCase.imageURL != "" {
	// 		imageURL := fmt.Sprintf(articleCase.imageURL, tickerSymbol)
	// 		linkURL := fmt.Sprintf(articleCase.linkURL, tickerSymbol)
	// 		photo := &tb.Photo{
	// 			File: tb.FromURL(imageURL),
	// 			Caption: fmt.Sprintf(
	// 				`\#%s [%s](%s)`,
	// 				tickerSymbol,
	// 				escape(articleCase.name),
	// 				linkURL,
	// 			),
	// 		}
	// 		b.Send(
	// 			to,
	// 			photo,
	// 			&tb.SendOptions{
	// 				ParseMode: tb.ModeMarkdownV2,
	// 			},
	// 		)
	// 	}
	// 	// if (len(commands) == 0 || contains(commands, "ark")) && contains(ARKTickets, ticketName) {
	// 	// 	log.Println("OK")
	// 	// 	log.Println("====")
	// 	// 	b.Send(
	// 	// 		to,
	// 	// 		fmt.Sprintf(
	// 	// 			"\\#%s [ARK](https://cathiesark.com/ark-combined-holdings-of-%s)",
	// 	// 			ticketName,
	// 	// 			strings.ToLower(ticketName),
	// 	// 		),
	// 	// 		&tb.SendOptions{
	// 	// 			ParseMode: tb.ModeMarkdownV2,
	// 	// 		},
	// 	// 	)
	// 	// }
	// })
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
