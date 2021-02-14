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
	InitializeArticleCases()
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
		// ownerID   = os.Getenv("OWNER_ID")   // you must add it to your config vars
		chatID = os.Getenv("CHAT_ID") // you must add it to your config vars
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
		url := fmt.Sprintf("https://tradingview.com/symbols/%s", ticker.symbol)
		result := &tb.ArticleResult{
			Title:       ticker.symbol,
			Description: ticker.description,
			HideURL:     true,
			URL:         url,
			ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticker.symbol),
		}
		result.SetContent(&tb.InputTextMessageContent{
			Text: fmt.Sprintf(`\#%s \- [%s](%s)`,
				ticker.symbol,
				escape(ticker.description),
				url,
			),
			ParseMode:      tb.ModeMarkdownV2,
			DisablePreview: true,
		})
		result.SetResultID("")
		results[0] = result
		for name, url := range ArticleCases {
			url := fmt.Sprintf(url, ticker.symbol)
			result := &tb.ArticleResult{
				Title:       name,
				Description: ticker.symbol,
				HideURL:     true,
				URL:         url,
			}
			result.SetContent(&tb.InputTextMessageContent{
				Text: fmt.Sprintf(`\#%s [%s](%s)`,
					ticker.symbol,
					escape(name),
					url,
				),
				ParseMode:      tb.ModeMarkdownV2,
				DisablePreview: true,
			})
			result.SetResultID(ticker.symbol + "=" + name)
			// results[i+1] = result
			results = append(results, result)
		}
		err = b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
			// SwitchPMText:      "SwitchPMText",
			// SwitchPMParameter: "SwitchPMParameter",
		})
		if err != nil {
			log.Println(err)
		}
	})
	b.Handle(tb.OnChosenInlineResult, func(r *tb.ChosenInlineResult) {
		// incoming inline queries
		log.Println("====")
		log.Println(r.MessageID)
		log.Println(r.ResultID)
		log.Println(r.Query)
		log.Println(r.From.ID)
		log.Println("====")
		if r.ResultID == "" {
			return
		}
		resultID := strings.Split(r.ResultID, "=")
		tickerSymbol := resultID[0]
		name := resultID[1]
		log.Println(name)
		log.Println(tickerSymbol)
		// ticketName := r.ResultID
		// TODO: to
		to := tb.ChatID(parseInt(chatID))
		// commands := make([]string, 0)
		// for _, param := range strings.Split(r.Query, " ") {
		// 	if strings.HasPrefix(param, "#") || strings.HasPrefix(param, "$") {
		// 		continue
		// 	}
		// 	commands = append(commands, param)
		// }
		if name == "finviz.com" {
			url := ArticleCases[name]
			screenshot := Screenshot(url)
			photo := &tb.Photo{
				File: tb.FromReader(bytes.NewReader(screenshot)),
				Caption: fmt.Sprintf(
					`\#%s [%s](%s)`,
					tickerSymbol,
					escape(name),
					url,
				),
			}
			b.Send(
				to,
				photo,
				&tb.SendOptions{
					ParseMode: tb.ModeMarkdownV2,
				},
			)
		}
		if name == "stockscores.com" {
			url := ArticleCases[name]
			photo := &tb.Photo{
				File: tb.FromURL(url),
				Caption: fmt.Sprintf(
					`\#%s [%s](%s)`,
					tickerSymbol,
					escape(name),
					url,
				),
			}
			b.Send(
				to,
				photo,
				&tb.SendOptions{
					ParseMode: tb.ModeMarkdownV2,
				},
			)

		}
		// if (len(commands) == 0 || contains(commands, "ark")) && contains(ARKTickets, ticketName) {
		// 	log.Println("OK")
		// 	log.Println("====")
		// 	b.Send(
		// 		to,
		// 		fmt.Sprintf(
		// 			"\\#%s [ARK](https://cathiesark.com/ark-combined-holdings-of-%s)",
		// 			ticketName,
		// 			strings.ToLower(ticketName),
		// 		),
		// 		&tb.SendOptions{
		// 			ParseMode: tb.ModeMarkdownV2,
		// 		},
		// 	)
		// }
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
