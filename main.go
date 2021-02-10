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
	b.Handle(tb.OnText, func(m *tb.Message) {
		log.Println("+++++")
		log.Println(m.ID)
		log.Println(m.InlineID)
		log.Println(m.Sender.ID)
		log.Println(m.Via.IsBot)
		log.Println(m.Via.Username)
		log.Println("+++++")
		// if m.Via.ID != b.Me.ID {
		// 	return
		// }
		// err := b.Delete(
		// 	&tb.StoredMessage{
		// 		MessageID: strconv.Itoa(m.ID),
		// 		ChatID:    parseInt(chatID),
		// 	},
		// )
		// if err != nil {
		// 	log.Println(err)
		// }
	})
	b.Handle(tb.OnQuery, func(q *tb.Query) {
		log.Println("*****")
		log.Println(q.Text)
		log.Println(q.ID)
		log.Println(q.From.ID)
		log.Println(ownerID)
		log.Println("*****")
		// TODO: разрешить всем админам чата
		// chat, err := b.ChatByID(chatID)
		// if err != nil {
		// 	log.Println(err)
		// }
		// chat
		// if strconv.Itoa(q.From.ID) != ownerID {
		// 	return
		// }
		var search string
		for _, param := range strings.Split(q.Text, " ") {
			if strings.HasPrefix(param, "#") {
				search = strings.TrimLeft(param, "#")
				break
			}
			if strings.HasPrefix(param, "$") {
				search = strings.TrimLeft(param, "$")
				break
			}
		}
		tickets := GetTickets(search)
		results := make(tb.Results, len(tickets)) // []tb.Result
		for i, ticket := range tickets {
			url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/", ticket.name)
			result := &tb.ArticleResult{
				Title:       ticket.name,
				Description: ticket.description,
				Text:        "OK",
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
		log.Println("=====")
		log.Println(r.MessageID)
		log.Println(r.ResultID)
		log.Println(r.Query)
		log.Println(r.From.ID)
		log.Println("=====")
		ticketName := r.ResultID
		to := tb.ChatID(parseInt(chatID))
		commands := make([]string, 0)
		for _, param := range strings.Split(r.Query, " ") {
			if strings.HasPrefix(param, "#") || strings.HasPrefix(param, "$") {
				continue
			}
			commands = append(commands, param)
		}
		if len(commands) == 0 || contains(commands, "finviz") {
			screenshot := Screenshot(ticketName)
			photo := &tb.Photo{
				File: tb.FromReader(bytes.NewReader(screenshot)),
				Caption: fmt.Sprintf(
					"\\#%[1]s [finviz](https://finviz.com/quote.ashx?t=%[1]s)",
					ticketName,
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
		if (len(commands) == 0 || contains(commands, "ark")) && contains(ARKTickets, ticketName) {
			b.Send(
				to,
				fmt.Sprintf(
					"\\#%s [ARK](https://cathiesark.com/ark-combined-holdings-of-%s)",
					ticketName,
					strings.ToLower(ticketName),
				),
				&tb.SendOptions{
					ParseMode: tb.ModeMarkdownV2,
				},
			)
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
