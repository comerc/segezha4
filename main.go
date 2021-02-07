package main

import (
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	// log.Fatal("err1234")

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

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!")
	})

	b.Handle("/text", func(m *tb.Message) {
		b.Send(m.Sender, "You entered "+m.Text)
	})

	b.Handle("/payload", func(m *tb.Message) {
		b.Send(m.Sender, "You entered "+m.Payload)
	})

	inlineBtn1 := tb.InlineButton{
		Unique: "moon",
		Text:   "Moon 🌚",
	}

	inlineBtn2 := tb.InlineButton{
		Unique: "sun",
		Text:   "Sun 🌞",
	}

	b.Handle(&inlineBtn1, func(c *tb.Callback) {
		// Required for proper work
		b.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		// Send messages here
		b.Send(c.Sender, "Moon says 'Hi'!")
	})

	b.Handle(&inlineBtn2, func(c *tb.Callback) {
		b.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		b.Send(c.Sender, "Sun says 'Hi'!")
	})

	inlineKeys := [][]tb.InlineButton{
		// []tb.InlineButton{inlineBtn1, inlineBtn2},
		{inlineBtn1, inlineBtn2},
	}

	b.Handle("/pick_time", func(m *tb.Message) {
		b.Send(
			m.Sender,
			"Day or night, you choose",
			&tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		// incoming inline queries
		log.Println(q.From)
		log.Println(q.Text)
	})

	b.Start()
}

// package main

// import (
// 	"io"
// 	"net/http"
// 	"os"
// )

// func hello(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "Hello World!")
// }

// func main() {
// 	port := os.Getenv("PORT")
// 	http.HandleFunc("/", hello)
// 	http.ListenAndServe(":"+port, nil)
// }
