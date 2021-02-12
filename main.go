package main

import (
	// "bytes"
	// "fmt"
	"log"
	"os"
	"strconv"

	// "strings"

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
		log.Println("****")
		log.Println(q.Text)
		log.Println(q.ID)
		log.Println(q.From.ID)
		log.Println("****")
		if q.Text == "" {
			return
		}
		results := make(tb.Results, len(ArticleCases)) // []tb.Result
		for i, articleCase := range ArticleCases {
			result := &tb.ArticleResult{
				Title:       articleCase.name,
				Text:        "OK",
				HideURL:     true,
				URL:         articleCase.url,
				ThumbURL:    "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAABACAYAAABcIPRGAAABVUlEQVRoQ+2YUQ6CQAxEyzk8rHpvjWaNGFEWi9NpxYw/JIaWeZ3uLmWwjf+Gjes3AVQ7KAfkAFiBv2ihC7kIBzPbg8/ohjcH2ADt4TSILAAaRCYABSIbIByiAiAUogogDKISIATCA7Azs5OZnTtXZBuGt1cPAOu0hsU3C6sAQsRXAYSJrwAIFZ8NEC4+E4AiPguAJj4DABV/vA8C3XmCuY1GiG85Fk9sFkCk+HEam83JAGCI70L8GkDr+bFtenPwpEAMgG/fMj3i35xgAayFWCN+AsEE8EJ8I/4BwQb4BIGIv0FkAPQgYPGZAK8QIeKzAUaI52tvq3T/n9VCbkEzNy7O3AJASuuMlQOszypOA5Y//2sNeMsI3Kc1oDUAtE8LVQttvoXADuCGe84BrgIwuwDAAsLhcgAuIZhADoAFhMOrDykByAG4B8AEcgAsIBwuB+ASggmuXaNxljt4uNQAAAAASUVORK5CYII=",
				ThumbWidth:  64,
				ThumbHeight: 64,
			}
			result.SetResultID(strconv.Itoa(i))
			results[i] = result
		}
		err = b.Answer(q, &tb.QueryResponse{
			Results:           results,
			CacheTime:         60, // a minute
			SwitchPMText:      "SwitchPMText",
			SwitchPMParameter: "SwitchPMParameter",
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
