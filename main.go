package main

import (
	"fmt"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token     = os.Getenv("TOKEN")      // you must add it to your config vars
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
		log.Println("****")
		tickets := GetTickets(q.Text)
		results := make(tb.Results, len(tickets)) // []tb.Result
		for i, ticket := range tickets {
			result := &tb.ArticleResult{
				Title:       ticket.name,
				Description: ticket.description,
				HideURL:     true,
				URL:         fmt.Sprintf("https://stockanalysis.com/stocks/%s/", ticket.name),
				ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", "GM"),
				// ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticket.name),
			}
			text := fmt.Sprintf("$%s \\- %s", ticket.name, ticket.description)
			// if contains(ARKTickets, ticket.name) {
			// 	text += fmt.Sprintf(" \\([ARK](https://cathiesark.com/ark-combined-holdings-of-%s)\\)", strings.ToLower(ticket.name))
			// }
			result.SetContent(&tb.InputTextMessageContent{
				Text:      text,
				ParseMode: tb.ModeMarkdownV2,
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
	// b.Handle(tb.OnChosenInlineResult, func(r *tb.ChosenInlineResult) {
	// 	// incoming inline queries
	// 	log.Println("====")
	// 	log.Println(r.MessageID)
	// 	log.Println(r.ResultID)
	// 	log.Println(r.Query)
	// 	log.Println(r.From.ID)
	// 	log.Println(r.From.Recipient())
	// 	log.Println("====")
	// 	i, err := strconv.ParseInt(chatID, 10, 64)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	to := tb.ChatID(i)
	// 	// photo := &tb.Photo{File: tb.FromURL("https://pp.vk.me/c627626/v627626512/2a627/7dlh4RRhd24.jpg")}
	// 	ticketName := r.ResultID
	// 	photo := &tb.Photo{
	// 		File:    tb.FromURL("https://firebasestorage.googleapis.com/v0/b/minsk8-2.appspot.com/o/8b98f59a-155b-464c-898f-1c04cfa86969.jpg?alt=media&token=2628e0bf-d11d-403f-98ac-b09fff126831"),
	// 		Caption: "#" + ticketName + " finviz",
	// 		// "https://finviz.com/quote.ashx?t=" + ticketName
	// 	}
	// 	b.Send(to, photo)
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
