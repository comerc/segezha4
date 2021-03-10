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

// TODO: отправлять через runBackgroundTask() информеры про фьючи
// TODO: использовать символы тикеров в качестве команд: /TSLA
// TODO: подключить ETF-ки https://etfdb.com/screener/
// TODO: выдавать сообщение sendLink, а по готовности основного ответа - его удалять
// TODO: кнопки под полем ввода в приватном чате для: inline mode, help, search & all,
// TODO: реализовать румтур
// TODO: поиск по ticker.title
// TODO: README
// TODO: svg to png
// TODO: добавить тайм-фрейм #BABA?15M
// TODO: добавить медленную скользящую #BABA?50EMA / 100EMA / 200EMA
// TODO: параллельная обрарботка https://gobyexample.ru/worker-pools.html
// TODO: добавить ETF, например ARKK
// TODO: добавить биток GBTC
// TODO: добавить опционы с investing.com
// TODO: не успевает загрузить картинку tipranks.com (показывает колёсики)

func main() {
	var (
		// port      = os.Getenv("PORT")
		// publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		chatID = os.Getenv("TELEBOT_CHAT_ID") // you must add it to your config vars
		token  = os.Getenv("TELEBOT_SECRET")  // you must add it to your config vars
	)
	// webhook := &tb.Webhook{
	// 	Listen:   ":" + port,
	// 	Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	// }
	pref := tb.Settings{
		// URL:    "https://api.bots.mn/telegram/",
		Token: token,
		// Poller: webhook,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
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
			linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
		log.Println("****")
		if m.Sender != nil {
			log.Println(m.Sender.Username)
			log.Println(m.Sender.FirstName)
			log.Println(m.Sender.LastName)
		}
		log.Println(m.Chat.Username)
		log.Println(m.Text)
		log.Println("****")
		if m.Text == "/ids" {
			sendFinvizIDs(b, m.Chat.ID)
		} else if m.Text == "/us" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabUS)
		} else if m.Text == "/europe" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabEurope)
		} else if m.Text == "/asia" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabAsia)
		} else if m.Text == "/fx" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabFX)
		} else if m.Text == "/rates" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabRates)
		} else if m.Text == "/futures" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabFutures)
		} else if m.Text == "/crypto" {
			sendMarketWatchIDs(b, m.Chat.ID, ss.MarketWatchTabCrypto)
		} else if m.Text == "/vix" {
			sendBarChart(b, m.Chat.ID, "$VIX")
		} else if m.Text == "/spy" {
			sendBarChart(b, m.Chat.ID, "SPY")
		} else if m.Text == "/index" {
			sendBarChart(b, m.Chat.ID, "$INX")
			sendBarChart(b, m.Chat.ID, "$NASX")
			sendBarChart(b, m.Chat.ID, "$DOWI")
		} else if m.Text == "/volume" {
			sendBarChart(b, m.Chat.ID, "SPY")
			sendBarChart(b, m.Chat.ID, "QQQ")
			sendBarChart(b, m.Chat.ID, "DOW")
		} else if m.Text == "/map" {
			sendFinvizMap(b, m.Chat.ID)
		} else if m.Text == "/fear" {
			sendFear(b, m.Chat.ID)
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
					if !result {
						sendError(b, m, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
						result = true
					}
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
				case "!":
					articleCase = GetExactArticleCase("finviz.com")
					result = sendScreenshotForFinviz(b, m, articleCase, ticker)
					// result = sendScreenshotForPage(b, m, articleCase, ticker)
					if !result {
						sendError(b, m, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
						result = true
					}
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
	go runBackgroundTask(b, int64(strToInt(chatID)))
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

func escapeURL(s string) string {
	re := regexp.MustCompile("[(|)]")
	return re.ReplaceAllString(s, `\$0`)
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
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

func runBackgroundTask(b *tb.Bot, chatID int64) {
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		h := t.UTC().Hour()
		m := t.Minute()
		const d = 30
		if h == 14 && m >= 30 || h > 14 && h < 21 || h == 21 && m < d {
			if m%d == 0 && t.Second() == 4 {
				sendFinvizIDs(b, chatID)
				sendFinvizMap(b, chatID)
				sendBarChart(b, chatID, "$VIX")
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabUS)
				if h >= 8 || h <= 17 {
					sendMarketWatchIDs(b, chatID, ss.MarketWatchTabEurope)
				}
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabRates)
			}
		} else if m == 0 {
			if h >= 0 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabFutures)
			}
			if h >= 8 || h <= 17 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabEurope)
			}
			if h >= 0 || h <= 8 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabAsia)
			}
			if h >= 0 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabRates)
			}
			// sendMarketWatchIDs(b, chatID, ss.MarketWatchTabFX)
			// sendMarketWatchIDs(b, chatID, ss.MarketWatchTabCrypto)
		}
	}
}

func sendBarChart(b *tb.Bot, chatID int64, symbol string) bool {
	volume, height, tag := func() (string, string, string) {
		if strings.HasPrefix(symbol, "$") {
			return "0", "625", ""
		}
		return "total", "500", `\#`
	}()
	linkURL := "https://www.barchart.com/stocks/quotes/%s/technical-chart%s?plot=CANDLE&volume=%s&data=I:5&density=L&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(50);EXPMA(20);EXPMA(200);WMA(9);EXPMA(500)&sym=%[1]s&grid=1&height=%[4]s&studyheight=200"
	screenshot := ss.MakeScreenshotForBarChart(fmt.Sprintf(linkURL, symbol, "/fullscreen", volume, height))
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			"%s[%s](%s)",
			escape(by(tag+symbol)),
			escape("barchart.com"),
			escapeURL(fmt.Sprintf(linkURL, symbol, "", volume, height)),
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

func sendFinvizMap(b *tb.Bot, chatID int64) bool {
	linkURL := "https://finviz.com/map.ashx?t=sec"
	screenshot := ss.MakeScreenshotForFinvizMap(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			"%s[%s](%s)",
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

func sendFear(b *tb.Bot, chatID int64) bool {
	linkURL := "https://money.cnn.com/data/fear-and-greed/"
	screenshot := ss.MakeScreenshotForFear(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			"%s[%s](%s)",
			escape(by("Fear & Greed Index")),
			escape("money.cnn.com"),
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

func sendFinvizIDs(b *tb.Bot, chatID int64) bool {
	linkURL := "https://finviz.com/"
	screenshot := ss.MakeScreenshotForFinvizIDs(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			"%s[%s](%s)",
			escape(by("IDs")),
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

func sendMarketWatchIDs(b *tb.Bot, chatID int64, tab ss.MarketWatchTab) bool {
	linkURL := "https://www.marketwatch.com/"
	screenshot := ss.MakeScreenshotForMarketWatchIDs(linkURL, tab)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			"%s[%s](%s)",
			escape(by("IDs")),
			escape("marketwatch.com"),
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
