package main

import (
	// "bytes"
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
		// ownerID   = os.Getenv("OWNER_ID")   // you must add it to your config vars
		// chatID    = os.Getenv("CHAT_ID")    // you must add it to your config vars
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
		url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/company/", ticker.symbol)
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
		for i, articleCase := range ArticleCases {
			url := fmt.Sprintf(articleCase.url, ticker.symbol)
			result := &tb.ArticleResult{
				Title:       articleCase.name,
				Description: ticker.symbol,
				HideURL:     true,
				URL:         url,
			}
			result.SetContent(&tb.InputTextMessageContent{
				Text: fmt.Sprintf(`\#%s by [%s](%s)`,
					ticker.symbol,
					escape(articleCase.name),
					url,
				),
				ParseMode:      tb.ModeMarkdownV2,
				DisablePreview: articleCase.hasPreview != true,
			})
			result.SetResultID(articleCase.name + "=" + ticker.symbol)
			results[i+1] = result
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
		articleCaseName := resultID[0]
		tickerSymbol := resultID[1]
		log.Println(articleCaseName)
		log.Println(tickerSymbol)
		// ticketName := r.ResultID
		// to := tb.ChatID(parseInt(chatID))
		// commands := make([]string, 0)
		// for _, param := range strings.Split(r.Query, " ") {
		// 	if strings.HasPrefix(param, "#") || strings.HasPrefix(param, "$") {
		// 		continue
		// 	}
		// 	commands = append(commands, param)
		// }
		// if len(commands) == 0 || contains(commands, "finviz") {
		// 	screenshot := Screenshot(ticketName)
		// 	photo := &tb.Photo{
		// 		File: tb.FromReader(bytes.NewReader(screenshot)),
		// 		Caption: fmt.Sprintf(
		// 			"\\#%[1]s [finviz](https://finviz.com/quote.ashx?t=%[1]s)",
		// 			ticketName,
		// 		),
		// 	}
		// 	b.Send(
		// 		to,
		// 		photo,
		// 		&tb.SendOptions{
		// 			ParseMode: tb.ModeMarkdownV2,
		// 		},
		// 	)
		// }
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
