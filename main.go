package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	ss "github.com/comerc/segezha4/screenshot"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: кнопки под полем ввода в приватном чате для: inline mode, help, search & all,
// TODO: bold для тикеров
// TODO: реализовать румтур
// TODO: поиск по ticker.title
// TODO: README
// TODO: svg to png
// TODO: добавить тайм-фрейм #BABA?15M
// TODO: добавить медленную скользящую #BABA?50EMA / 100EMA / 200EMA
// TODO: параллельная обрарботка https://gobyexample.ru/worker-pools.html
// TODO: добавить ETF, например ARKK
// TODO: добавить биток GBTC
// TODO: добавить https://stockcharts.com/
// TODO: не успевает загрузить картинку tipranks.com (показывает колёсики)

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
		// URL:    "https://api.bots.mn/telegram/",
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
		results := make(tb.Results, len(ArticleCases)) // []tb.Result
		for i, articleCase := range ArticleCases {
			linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
			var result *tb.ArticleResult
			if i == 0 {
				result = &tb.ArticleResult{
					Title:       fmt.Sprintf("%s #%s", articleCase.name, ticker.symbol),
					Description: ticker.title,
					HideURL:     true,
					URL:         linkURL,
					ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticker.symbol), // from stockanalysis.com
				}
			} else {
				title := fmt.Sprintf("%s #%s", articleCase.name, ticker.symbol)
				if articleCase.screenshotMode != "" {
					title += " 🎁"
				}
				result = &tb.ArticleResult{
					Title:       title,
					Description: articleCase.description,
					HideURL:     true,
					URL:         linkURL,
				}
			}
			result.SetContent(&tb.InputTextMessageContent{
				Text: fmt.Sprintf("/info %s %s",
					articleCase.name,
					ticker.symbol,
				),
				DisablePreview: true,
			})
			result.SetResultID(ticker.symbol + "=" + articleCase.name)
			results[i] = result
		}
		err = b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})
		if err != nil {
			log.Println(err)
		}
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		log.Println(m.Text)
		if m.Text == "/map" {
			sendFinvizMap(b, m.Chat.ID)
		} else if strings.HasPrefix(m.Text, "/info ") {
			re := regexp.MustCompile(",")
			payload := re.ReplaceAllString(m.Payload, " ")
			arguments := strings.Split(payload, " ")
			symbols := arguments[1:]
			if len(symbols) == 0 {
				sendError(b, m, "Empty symbols")
				return
			}
			articleCaseName := arguments[0]
			articleCase := GetExactArticleCase(articleCaseName)
			if articleCase == nil {
				sendError(b, m, "Invalid command")
				return
			}
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
					symbol = symbol[1:]
				}
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendError(b, m, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				var result bool
				switch articleCase.screenshotMode {
				case ScreenshotModePage:
					result = sendScreenshotForPage(b, m, articleCase, ticker)
				case ScreenshotModeImage:
					result = sendImage(b, m, articleCase, ticker)
					// result = sendScreenshotForImage(b, m, articleCase, ticker)
				case ScreenshotModeFinviz:
					result = sendScreenshotForFinviz(b, m, articleCase, ticker)
				case ScreenshotModeMarketWatch:
					result = sendScreenshotForMarketWatch(b, m, articleCase, ticker)
				case ScreenshotModeMarketBeat:
					result = sendScreenshotForMarketBeat(b, m, articleCase, ticker)
				case ScreenshotModeCathiesArk:
					result = sendScreenshotForCathiesArk(b, m, articleCase, ticker)
				default:
					result = false
				}
				if !result {
					sendLink(b, m, articleCase, ticker)
				}
			}
			// err := b.Delete(
			// 	&tb.StoredMessage{
			// 		MessageID: strconv.Itoa(m.ID),
			// 		ChatID:    m.Chat.ID,
			// 	},
			// )
			// if err != nil {
			// 	log.Println(err)
			// }
		} else {
			// simple command mode
			re := regexp.MustCompile(`(^|[ ])#([A-Za-z]+)(\?!|\?|!)`)
			matches := re.FindAllStringSubmatch(m.Text, -1)
			for _, match := range matches {
				symbol := match[2]
				mode := match[3]
				// log.Println(symbol + mode)
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendError(b, m, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				var (
					articleCase *ArticleCase
					result      bool
				)
				// TODO: var modes map[string]myFunc https://golangbot.com/first-class-functions/
				switch mode {
				case "?!":
					articleCase = GetExactArticleCase("marketwatch.com")
					result = sendScreenshotForMarketWatch(b, m, articleCase, ticker)
					// result = sendScreenshotForPage(b, m, articleCase, ticker)
					// articleCase = GetExactArticleCase("shortvolume.com")
					// result = sendImage(b, m, articleCase, ticker)
					// articleCase = GetExactArticleCase("shortvolume.com")
					// result = sendScreenshotForImage(b, m, articleCase, ticker)
				case "?":
					articleCase = GetExactArticleCase("stockscores.com")
					result = sendImage(b, m, articleCase, ticker)
					// articleCase = GetExactArticleCase("stockscores.com")
					// result = sendScreenshotForImage(b, m, articleCase, ticker)
				case "!":
					articleCase = GetExactArticleCase("finviz.com")
					// result = sendScreenshotForPage(b, m, articleCase, ticker)
					result = sendScreenshotForFinviz(b, m, articleCase, ticker)
				default:
					log.Println("Invalid simple command mode")
					result = true
				}
				if !result {
					sendLink(b, m, articleCase, ticker)
				}
			}
		}
	})
	b.Start()
	// go backgroundTask(b, int64(strToInt(chatID)))
	// // This print statement will be executed before
	// // the first `tock` prints in the console
	// log.Println("The rest of my application can continue")
	// // here we use an empty select{} in order to keep
	// // our main function alive indefinitely as it would
	// // complete before our backgroundTask has a chance
	// // to execute if we didn't.
	// select {}
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

// func deleteCommand(b *tb.Bot, m *tb.Message) {
// 	err := b.Delete(
// 		&tb.StoredMessage{
// 			MessageID: strconv.Itoa(m.ID),
// 			ChatID:    m.Chat.ID,
// 		},
// 	)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func sendInformer(b *tb.Bot, m *tb.Message, photo *tb.Photo) {
// 	_, err := b.Send(
// 		tb.ChatID(m.Chat.ID),
// 		photo,
// 		&tb.SendOptions{
// 			ParseMode: tb.ModeMarkdownV2,
// 		},
// 	)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func getUserLink(u *tb.User) string {
// 	if u.Username != "" {
// 		return fmt.Sprintf("@%s", u.Username)
// 	}
// 	return fmt.Sprintf("[%s](tg://user?id=%d)", u.FirstName, u.ID)
// }

func sendScreenshotForPage(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForPage(linkURL, articleCase.x, articleCase.y, articleCase.width, articleCase.height)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendScreenshotForFinviz(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForFinviz(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendScreenshotForMarketWatch(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForMarketWatch(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendScreenshotForMarketBeat(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForMarketBeat(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendScreenshotForCathiesArk(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForCathiesArk(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendScreenshotForImage(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol)
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	screenshot := ss.MakeScreenshotForImage(imageURL, articleCase.width, articleCase.height)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	screenshot = nil
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendImage(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) bool {
	imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol, time.Now().Unix())
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	photo := &tb.Photo{
		File: tb.FromURL(imageURL),
		Caption: fmt.Sprintf(
			`\#%s %s[%s](%s)`,
			ticker.symbol,
			escape(by(articleCase.description)),
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func sendLink(b *tb.Bot, m *tb.Message, articleCase *ArticleCase, ticker *Ticker) {
	description := func() string {
		if articleCase.name == ArticleCases[0].name {
			return ticker.title
		}
		return articleCase.description
	}()
	linkURL := fmt.Sprintf(articleCase.linkURL, ticker.symbol)
	text := fmt.Sprintf(`\#%s %s[%s](%s)`,
		ticker.symbol,
		escape(by(description)),
		escape(articleCase.name),
		linkURL,
		// getUserLink(m.Sender),
	)
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		text,
		&tb.SendOptions{
			ParseMode:             tb.ModeMarkdownV2,
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func sendError(b *tb.Bot, m *tb.Message, text string) {
	_, err := b.Send(
		tb.ChatID(m.Chat.ID),
		text,
		&tb.SendOptions{
			ParseMode:             tb.ModeMarkdownV2,
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func by(s string) string {
	if s == "" {
		return "by "
	}
	return s + " by "
}

// func backgroundTask(b *tb.Bot, chatID int64) {
// 	ticker := time.NewTicker(1 * time.Second)
// 	for t := range ticker.C {
// 		log.Println("Tick at", t.Minute(), t.Minute()%10, t.Second())
// 		// t.Minute()%10 == 0 &&
// 		if t.Second() == 4 {
// 			if sendFinvizMap(b, chatID) {
// 				log.Println("Send map")
// 			}
// 		}
// 	}
// }

func sendFinvizMap(b *tb.Bot, chatID int64) bool {
	linkURL := "https://finviz.com/map.ashx?t=sec"
	screenshot := ss.MakeScreenshotForFinvizMap(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`%s[%s](%s)`,
			escape(by("Map")),
			escape("finviz.com"),
			linkURL,
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
		photo,
		&tb.SendOptions{
			ParseMode: tb.ModeMarkdownV2,
		},
	)
	photo = nil
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func strToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println(s, "is not an integer.")
		return 0
	}
	return n
}
