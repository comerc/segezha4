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
