package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
		ownerID   = os.Getenv("OWNER_ID")   // you must add it to your config vars
		chatID    = os.Getenv("CHAT_ID")    // you must add it to your config vars
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
	// b.Handle("/hello", func(m *tb.Message) {
	// 	b.Send(m.Sender, "Hi!")
	// })
	// b.Handle("/text", func(m *tb.Message) {
	// 	b.Send(m.Sender, "You entered "+m.Text)
	// })
	// b.Handle("/payload", func(m *tb.Message) {
	// 	b.Send(m.Sender, "You entered "+m.Payload)
	// })
	// inlineBtn1 := tb.InlineButton{
	// 	Unique: "moon",
	// 	Text:   "Moon 🌚",
	// }
	// inlineBtn2 := tb.InlineButton{
	// 	Unique: "sun",
	// 	Text:   "Sun 🌞",
	// }
	// b.Handle(&inlineBtn1, func(c *tb.Callback) {
	// 	// Required for proper work
	// 	b.Respond(c, &tb.CallbackResponse{
	// 		ShowAlert: false,
	// 	})
	// 	// Send messages here
	// 	b.Send(c.Sender, "Moon says 'Hi'!")
	// })
	// b.Handle(&inlineBtn2, func(c *tb.Callback) {
	// 	b.Respond(c, &tb.CallbackResponse{
	// 		ShowAlert: false,
	// 	})
	// 	b.Send(c.Sender, "Sun says 'Hi'!")
	// })
	// inlineKeys := [][]tb.InlineButton{
	// 	// []tb.InlineButton{inlineBtn1, inlineBtn2},
	// 	{inlineBtn1, inlineBtn2},
	// }
	// b.Handle("/pick_time", func(m *tb.Message) {
	// 	b.Send(
	// 		m.Sender,
	// 		"Day or night, you choose",
	// 		&tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	// })
	b.Handle(tb.OnQuery, func(q *tb.Query) {
		log.Println("****")
		log.Println(q.Text)
		log.Println(len(q.Text))
		log.Println(q.ID)
		log.Println(q.From.ID)
		log.Println(int64(q.From.ID) == parseInt(ownerID))
		log.Println("OK!!!!")
		log.Println("****")
		// TODO: разрешить всем админам чата
		// chat, err := b.ChatByID(chatID)
		// if err != nil {
		// 	log.Println(err)
		// }
		// chat
		tickets := GetTickets(q.Text)
		results := make(tb.Results, len(tickets)) // []tb.Result
		for i, ticket := range tickets {
			url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/", ticket.name)
			result := &tb.ArticleResult{
				Title:       ticket.name,
				Description: ticket.description,
				HideURL:     true,
				URL:         url,
				ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticket.name),
			}
			result.SetContent(&tb.InputTextMessageContent{
				Text: fmt.Sprintf("$%s \\- [%s](%s)",
					ticket.name,
					strings.Replace(ticket.description, ".", "\\.", -1),
					url,
				),
				ParseMode:      tb.ModeMarkdownV2,
				DisablePreview: true,
			})
			// result.SetReplyMarkup(inlineKeys)
			// needed to set a unique string ID for each result
			result.SetResultID(ticket.name)
			results[i] = result
			// TODO: max 50
		}
		err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
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
		log.Println(r.From.Recipient())
		log.Println("====")
		to := tb.ChatID(parseInt(chatID))
		ticketName := r.ResultID
		screenshot := Screenshot(ticketName)
		photo := &tb.Photo{
			File: tb.FromReader(bytes.NewReader(screenshot)),
			// FromURL("https://firebasestorage.googleapis.com/v0/b/minsk8-2.appspot.com/o/8b98f59a-155b-464c-898f-1c04cfa86969.jpg?alt=media&token=2628e0bf-d11d-403f-98ac-b09fff126831"),
			// Caption: "#" + ticketName + " finviz",
			// "https://finviz.com/quote.ashx?t=" + ticketName
			Caption: fmt.Sprintf(
				"\\#%[1]s [finviz](https://finviz.com/quote.ashx?t=%[1]s)",
				ticketName,
			),
			ParseMode: tb.ModeMarkdownV2,
		}
		b.Send(to, photo)

		// if contains(ARKTickets, ticketName) {
		// 	b.Send(to,
		// 		// "\\#TSLA [ARK](https://cathiesark.com/ark-combined-holdings-of-tsla)",
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
