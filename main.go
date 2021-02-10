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
		log.Println("****")
		log.Println(q.Text)
		log.Println(len(q.Text))
		log.Println(q.ID)
		log.Println(q.From.ID)
		log.Println("****")

		// plusIcon := "https://pp.vk.me/c627626/v627626512/2a627/7dlh4RRhd24.jpg"
		// minusIcon := "https://pp.vk.me/c627626/v627626512/2a635/ILYe7N2n8Zo.jpg"
		// divideIcon := "https://pp.vk.me/c627626/v627626512/2a620/oAvUk7Awps0.jpg"
		// multiplyIcon := "https://pp.vk.me/c627626/v627626512/2a62e/xqnPMigaP5c.jpg"
		// errorIcon := "https://pp.vk.me/c627626/v627626512/2a67a/ZvTeGq6Mf88.jpg"
		// tslaIcon := "https://financemarker.ru/fa/fa_logos/TSLA.png"
		// nvdaIcon := "https://financemarker.ru/fa/fa_logos/NVDA.png"
		// vrtxIcon := "https://financemarker.ru/fa/fa_logos/VRTX.png"
		// twtrIcon := "https://financemarker.ru/fa/fa_logos/TWTR.png"
		// https://storage.googleapis.com/iexcloud-hl37opg/api/logos/TWTR.png

		// urls := []string{
		// 	plusIcon,
		// 	minusIcon,
		// 	divideIcon,
		// 	multiplyIcon,
		// 	errorIcon,
		// 	tslaIcon,
		// 	nvdaIcon,
		// 	vrtxIcon,
		// 	twtrIcon,
		// }

		tickets := GetTickets(q.Text)

		results := make(tb.Results, len(tickets)) // []tb.Result
		for i, ticket := range tickets {
			result := &tb.ArticleResult{
				Title:       ticket.name,        // "Title" + fmt.Sprint(i) + " *Bold*",
				Description: ticket.description, // "Description" + fmt.Sprint(i) + " *Bold*",
				Text:        "OK " + ticket.description,

				// URL:       "https://finviz.com/quote.ashx?t=LMT",
				// MIME:      "text/html",

				ThumbURL: "https://storage.googleapis.com/iexcloud-hl37opg/api/logos/" + ticket.name + ".png",
			}
			//  .PhotoResult{
			// 	URL: url,

			// 	// required for photos
			// 	ThumbURL: url,
			// }

			// result.SetContent(&tb.InputTextMessageContent{
			// 	Text: "Text" + fmt.Sprint(i) + " *Bold* [src](https://itsallwidgets.com/screenshots/app-2041.png)",

			// 	ParseMode: tb.ModeMarkdownV2,
			// })

			result.SetResultID(ticket.name)

			// result.SetResultID("TSLA" + fmt.Sprint(i))

			// result.SetReplyMarkup(inlineKeys)

			results[i] = result
			// needed to set a unique string ID for each result
			// results[i].SetResultID(strconv.Itoa(i))

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

		group := tb.ChatID(-1001439193008)

		// photo := &tb.Photo{File: tb.FromURL("https://pp.vk.me/c627626/v627626512/2a627/7dlh4RRhd24.jpg")}

		photo := &tb.Photo{File: tb.FromURL("https://firebasestorage.googleapis.com/v0/b/minsk8-2.appspot.com/o/8b98f59a-155b-464c-898f-1c04cfa86969.jpg?alt=media&token=2628e0bf-d11d-403f-98ac-b09fff126831")}

		b.Send(group, photo)
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
