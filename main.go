package main

import (
	"fmt"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
)

// type ArticleResultExt struct {
// 	tb.ArticleResult

// 	// Optional. Send Markdown or HTML, if you want Telegram apps to show
// 	// bold, italic, fixed-width text or inline URLs in the media caption.
// 	ParseMode tb.ParseMode `json:"parse_mode,omitempty"`
// }

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
		log.Println("****")
		log.Println(q.Text)
		log.Println("****")

		plusIcon := "https://pp.vk.me/c627626/v627626512/2a627/7dlh4RRhd24.jpg"
		minusIcon := "https://pp.vk.me/c627626/v627626512/2a635/ILYe7N2n8Zo.jpg"
		divideIcon := "https://pp.vk.me/c627626/v627626512/2a620/oAvUk7Awps0.jpg"
		multiplyIcon := "https://pp.vk.me/c627626/v627626512/2a62e/xqnPMigaP5c.jpg"
		errorIcon := "https://pp.vk.me/c627626/v627626512/2a67a/ZvTeGq6Mf88.jpg"

		urls := []string{
			plusIcon,
			minusIcon,
			divideIcon,
			multiplyIcon,
			errorIcon,
		}

		results := make(tb.Results, len(urls)) // []tb.Result
		for i, url := range urls {
			result := &tb.ArticleResult{
				Title:       "Title" + fmt.Sprint(i) + " *Bold*",
				Description: "Description" + fmt.Sprint(i) + " *Bold*",

				// URL:       "https://finviz.com/quote.ashx?t=LMT",
				// MIME:      "text/html",

				ThumbURL: url,
			}
			//  .PhotoResult{
			// 	URL: url,

			// 	// required for photos
			// 	ThumbURL: url,
			// }

			result.SetContent(&tb.InputTextMessageContent{
				Text: "Text" + fmt.Sprint(i) + " *Bold* [src](https://itsallwidgets.com/screenshots/app-2041.png)",

				ParseMode: tb.ModeMarkdownV2,
			})

			result.SetResultID("TSLA" + fmt.Sprint(i))

			result.SetReplyMarkup(inlineKeys)

			results[i] = result
			// needed to set a unique string ID for each result
			// results[i].SetResultID(strconv.Itoa(i))
		}

		err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})

		if err != nil {
			log.Println(err)
		}
	})

	// b.Handle(tb.OnQuery, func(q *tb.Query) {
	// 	// incoming inline queries
	// 	log.Println(q.From.Username)
	// 	log.Println(q.Text)
	// 	// err := b.Answer(q, &tb.QueryResponse{
	// 	// 	Results:   results,
	// 	// 	CacheTime: 60, // a minute
	// 	// })
	// 	// if err != nil {
	// 	// 	log.Println(err)
	// 	// }
	// 	// tb.PhotoResult

	// })

	b.Handle(tb.OnChosenInlineResult, func(r *tb.ChosenInlineResult) {
		// incoming inline queries
		log.Println("====")
		log.Println(r.MessageID)
		log.Println(r.ResultID)
		log.Println(r.Query)
		log.Println("====")

		// photo := &tb.Photo{File: tb.FromURL("https://pp.vk.me/c627626/v627626512/2a627/7dlh4RRhd24.jpg")}

		// log.Println(q.Text)
		// err := b.Answer(q, &tb.QueryResponse{
		// 	Results:   results,
		// 	CacheTime: 60, // a minute
		// })
		// if err != nil {
		// 	log.Println(err)
		// }
		// tb.PhotoResult

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
