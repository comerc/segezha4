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

// TODO: ARK - перемножать кол-во купленных и проданных акций
// TODO: добавить опционы с investing.com
// TODO: использовать символы тикеров в качестве команд: /TSLA (но #TSLA! тоже оставить, иначе потеряю возможность вставлять внутри текста)
// TODO: подключить ETF-ки, например ARKK https://etfdb.com/screener/
// TODO: выдавать сообщение sendLink, а по готовности основного ответа - его удалять
// TODO: кнопки под полем ввода в приватном чате для: inline mode, help, search & all,
// TODO: поиск по ticker.title
// TODO: README
// TODO: svg to png
// TODO: добавить тайм-фрейм #BABA?15M
// TODO: добавить медленную скользящую #BABA?50EMA / 100EMA / 200EMA
// TODO: параллельная обработка https://gobyexample.ru/worker-pools.html
// TODO: добавить биток GBTC
// TODO: выборка с графиками https://finviz.com/screener.ashx?v=212&t=ZM,BA,MU,MS,GE,AA
// TODO: https://stockcharts.com/h-sc/ui?s=$CPCE https://school.stockcharts.com/doku.php?id=market_indicators:put_call_ratio

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
	// b.Handle(tb.OnCallback, func(*tb.Callback) {
	// 	log.Println("OnCallback")
	// })
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
	messageHandler := func(m *tb.Message) {
		log.Println("****")
		if m.Sender != nil {
			log.Println(m.Sender.Username)
			log.Println(m.Sender.FirstName)
			log.Println(m.Sender.LastName)
		}
		log.Println(m.Chat.Username)
		var text string
		if m.Photo != nil {
			text = m.Caption
		} else {
			text = m.Text
		}
		log.Println(text)
		log.Println("****")
		for tab := range ss.MarketWatchTabs {
			if text == "/"+tab {
				sendMarketWatchIDs(b, m.Chat.ID, tab)
				return
			}
		}
		if text == "/start" || text == "/help" {
			help := `*Commands:*
/help - this message
/bb - Bull Or Bear
/map - S&P 500 1 Day Performance Map
/fear - Fear & Greed Index
/us - US Indexes
/europe - Europe Indexes
/asia - Asia Indexes
/fx - Currencies
/rates - Bonds
/futures - Futures
/crypto - Crypto Currencies
/vix - $VIX (15M)
/spy - SPY (15M)
/index - Indexes (15M): $INX, $NASX, $DOWI
/volume - Volumes (15M): SPY, QQQ, DOW
/finviz - batch mode, for example: /finviz TSLA ETSY
/info - batch mode, for example: /info finviz.com GS MS

*Inline Menu Mode:*
@TickerInfoBot TSLA

*Simple (Batch) Mode:*
#TSLA! #TSLA? #TSLA?? #TSLA?! #TSLA!!
`
			sendText(b, m.Chat.ID, escape(help))
		} else if text == "/bb" {
			sendFinvizBB(b, m.Chat.ID)
		} else if text == "/vix" {
			sendBarChart(b, m.Chat.ID, "$VIX")
		} else if text == "/spy" {
			sendBarChart(b, m.Chat.ID, "SPY")
		} else if text == "/index" {
			sendBarChart(b, m.Chat.ID, "$INX")
			sendBarChart(b, m.Chat.ID, "$NASX")
			sendBarChart(b, m.Chat.ID, "$DOWI")
		} else if text == "/volume" {
			sendBarChart(b, m.Chat.ID, "SPY")
			sendBarChart(b, m.Chat.ID, "QQQ")
			sendBarChart(b, m.Chat.ID, "DOW")
		} else if text == "/map" {
			sendFinvizMap(b, m.Chat.ID)
		} else if text == "/fear" {
			sendFear(b, m.Chat.ID)
		} else if strings.HasPrefix(text, "/finviz ") {
			re := regexp.MustCompile(",|[ ]+")
			payload := re.ReplaceAllString(strings.Trim(m.Payload, " "), " ")
			symbols := strings.Split(payload, " ")
			if len(symbols) == 0 {
				sendText(b, m.Chat.ID, "No symbols")
				return
			}
			executed := make([]string, 0)
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
					symbol = symbol[1:]
				}
				if Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				result := sendFinvizImage(b, m.Chat.ID, symbol)
				if !result {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
				}
			}
		} else if strings.HasPrefix(text, "/info ") {
			re := regexp.MustCompile(",|[ ]+")
			payload := re.ReplaceAllString(strings.Trim(m.Payload, " "), " ")
			arguments := strings.Split(payload, " ")
			symbols := arguments[1:]
			if len(symbols) == 0 {
				sendText(b, m.Chat.ID, "No symbols")
				return
			}
			articleCaseName := arguments[0]
			articleCase := GetExactArticleCase(articleCaseName)
			if articleCase == nil {
				sendText(b, m.Chat.ID, "Invalid command")
				return
			}
			executed := make([]string, 0)
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
					symbol = symbol[1:]
				}
				if Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				var result bool
				switch articleCase.screenshotMode {
				case ScreenshotModePage:
					result = sendScreenshotForPage(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeImage:
					result = sendImage(b, m.Chat.ID, articleCase, ticker)
					// result = sendScreenshotForImage(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeFinviz:
					result = sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
						result = true
					}
				case ScreenshotModeMarketWatch:
					result = sendScreenshotForMarketWatch(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeMarketBeat:
					result = sendScreenshotForMarketBeat(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeCathiesArk:
					result = sendScreenshotForCathiesArk(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeGuruFocus:
					result = sendScreenshotForGuruFocus(b, m.Chat.ID, articleCase, ticker)
				case ScreenshotModeTipRanks:
					result = sendScreenshotForTipRanks(b, m.Chat.ID, articleCase, ticker)
				default:
					result = false
				}
				if !result {
					sendLink(b, m.Chat.ID, articleCase, ticker)
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

			// } else if isEarnings(text) {
			// 	re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			// 	matches := re.FindAllStringSubmatch(text, -1)
			// 	if len(matches) == 0 {
			// 		return
			// 	}
			// 	symbol := matches[0][2]
			// 	ticker := GetExactTicker(symbol)
			// 	if ticker == nil {
			// 		sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
			// 		return
			// 	}
			// 	articleCase := GetExactArticleCase("marketwatch.com")
			// 	result := sendScreenshotForMarketWatch(b, m.Chat.ID, articleCase, ticker)
			// 	if !result {
			// 		sendLink(b, m.Chat.ID, articleCase, ticker)
			// 	}
		} else if isEarnings(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				if Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				articleCase := GetExactArticleCase("marketwatch.com")
				result := sendScreenshotForMarketWatch(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendLink(b, m.Chat.ID, articleCase, ticker)
				}
			}
		} else if isARK(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			executed = append(executed, "ARK")
			for _, match := range matches {
				symbol := match[2]
				if Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				articleCase := GetExactArticleCase("finviz.com")
				result := sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
					result = true
				}
			}
		} else if isIdeas(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				if Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
					continue
				}
				articleCase := GetExactArticleCase("finviz.com")
				result := sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
					result = true
				}
			}
		} else {
			// simple command mode
			// TODO: "#ZM!!"
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)(\?!|\?\?|\?|!!|!)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				mode := match[3]
				if Contains(executed, strings.ToUpper(symbol)+mode) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol)+mode)
				// log.Println(symbol + mode)
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)))
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
					result = sendScreenshotForMarketWatch(b, m.Chat.ID, articleCase, ticker)
					// result = sendScreenshotForPage(b, m.Chat.ID, articleCase, ticker)
					// articleCase = GetExactArticleCase("shortvolume.com")
					// result = sendImage(b, m.Chat.ID, articleCase, ticker)
					// articleCase = GetExactArticleCase("shortvolume.com")
					// result = sendScreenshotForImage(b, m.Chat.ID, articleCase, ticker)
				case "??":
					articleCase = GetExactArticleCase("barchart.com") // для sendLink
					result = sendBarChart(b, m.Chat.ID, ticker.symbol)
				case "?":
					articleCase = GetExactArticleCase("stockscores.com")
					result = sendImage(b, m.Chat.ID, articleCase, ticker)
				case "!!":
					articleCase = GetExactArticleCase("finviz.com")
					result = sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
						result = true
					}
					if !result {
						sendLink(b, m.Chat.ID, articleCase, ticker)
					}
					articleCase = GetExactArticleCase("gurufocus.com")
					result = sendScreenshotForGuruFocus(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendLink(b, m.Chat.ID, articleCase, ticker)
					}
					articleCase = GetExactArticleCase("marketbeat.com")
					result = sendScreenshotForMarketBeat(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendLink(b, m.Chat.ID, articleCase, ticker)
					}
					articleCase = GetExactArticleCase("tipranks.com")
					result = sendScreenshotForTipRanks(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendLink(b, m.Chat.ID, articleCase, ticker)
					}
					result = true
				case "!":
					articleCase = GetExactArticleCase("finviz.com")
					result = sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
					if !result {
						sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)))
						result = true
					}
				default:
					log.Println("Invalid simple command mode")
					result = true
				}
				if !result {
					sendLink(b, m.Chat.ID, articleCase, ticker)
				}
			}
		}

	}
	b.Handle(tb.OnText, messageHandler)
	b.Handle(tb.OnPhoto, messageHandler)
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
	re := regexp.MustCompile(`[.|\-|(|)|#|!]`)
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

// func sendInformer(b *tb.Bot, chatID int64, photo *tb.Photo) {
// 	_, err := b.Send(
// 		tb.ChatID(chatID),
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

func sendScreenshotForPage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForPage(linkURL, articleCase.x, articleCase.y, articleCase.width, articleCase.height)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForFinviz(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForFinviz(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForMarketWatch(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForMarketWatch(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForMarketBeat(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForMarketBeat(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForCathiesArk(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForCathiesArk(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForGuruFocus(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForGuruFocus(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForTipRanks(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForTipRanks(strings.ToLower(ticker.symbol))
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendScreenshotForImage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol)
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	screenshot := ss.MakeScreenshotForImage(imageURL, articleCase.width, articleCase.height)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
		),
	}
	_, err := b.Send(
		tb.ChatID(chatID),
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

func sendFinvizImage(b *tb.Bot, chatID int64, symbol string) bool {
	imageURL := fmt.Sprintf("https://charts2.finviz.com/chart.ashx?t=%s&ta=1&p=d&r=%d", strings.ToLower(symbol), time.Now().Unix())
	linkURL := fmt.Sprintf("https://finviz.com/quote.ashx?t=%s", strings.ToLower(symbol))
	photo := &tb.Photo{
		File: tb.FromURL(imageURL),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			strings.ToUpper(symbol),
			escape("finviz.com"),
			linkURL,
			// getUserLink(m.Sender),
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

func sendImage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
	imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol, time.Now().Unix())
	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
	photo := &tb.Photo{
		File: tb.FromURL(imageURL),
		Caption: fmt.Sprintf(
			`\#%s by [%s](%s)`,
			ticker.symbol,
			escape(articleCase.name),
			linkURL,
			// getUserLink(m.Sender),
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

func sendLink(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) {
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
		tb.ChatID(chatID),
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

func sendText(b *tb.Bot, chatID int64, text string) {
	_, err := b.Send(
		tb.ChatID(chatID),
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
		w := t.Weekday()
		if w == 6 || w == 0 {
			continue
		}
		h := t.UTC().Hour()
		m := t.Minute()
		s := t.Second()
		const (
			d      = 30
			summer = 1
		)
		if h == 14-summer && m >= 30 || h > 14-summer && h < 21-summer || h == 21-summer && m < d {
			if m%d == 0 && s == 15 {
				if h == 14-summer && m >= 30 {
					sendFear(b, chatID)
				}
				if h >= 15-summer {
					sendFinvizBB(b, chatID)
					sendFinvizMap(b, chatID)
				}
				sendBarChart(b, chatID, "$VIX")
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabUS)
				if h >= 8 && h <= 17 {
					sendMarketWatchIDs(b, chatID, ss.MarketWatchTabEurope)
				}
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabRates)
				// sendMarketWatchIDs(b, chatID, ss.MarketWatchTabFutures)
			}
		} else if m == 0 && s == 15 {
			if h >= 8 && h <= 17 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabEurope)
			}
			// SPB работает с 7 утра (MSK)
			if h >= 4 && h <= 9 {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabAsia)
			}
			if h >= 4 && h <= 14-summer {
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabRates)
				sendMarketWatchIDs(b, chatID, ss.MarketWatchTabFutures)
			}
			// sendMarketWatchIDs(b, chatID, ss.MarketWatchTabFX)
			// sendMarketWatchIDs(b, chatID, ss.MarketWatchTabCrypto)
		}
	}
}

func sendBarChart(b *tb.Bot, chatID int64, symbol string) bool {
	volume, height, tag := func() (string, string, string) {
		if strings.HasPrefix(symbol, "$") {
			return "0", "O", ""
		}
		return "total", "H", "#"
	}()
	linkURL := "https://www.barchart.com/stocks/quotes/%s/technical-chart%s?plot=CANDLE&volume=%s&data=I:15&density=%[4]s&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(50);EXPMA(20);EXPMA(200);WMA(9);EXPMA(500)&sym=%[1]s&grid=1&height=500&studyheight=200"
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
			`\#map by [%s](%s)`,
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
			`\#fear by [%s](%s)`,
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

func sendFinvizBB(b *tb.Bot, chatID int64) bool {
	linkURL := "https://finviz.com/"
	screenshot := ss.MakeScreenshotForFinvizBB(linkURL)
	if len(screenshot) == 0 {
		return false
	}
	photo := &tb.Photo{
		File: tb.FromReader(bytes.NewReader(screenshot)),
		Caption: fmt.Sprintf(
			`\#bb Bull or Bear by [%s](%s)`,
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
			`\#%s by [%s](%s)`,
			escape(tab),
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

func isEarnings(text string) bool {
	re := regexp.MustCompile("#ОТЧЕТ")
	return re.FindStringIndex(text) != nil
}

func isARK(text string) bool {
	re := regexp.MustCompile("#ARK")
	return re.FindStringIndex(text) != nil
}

func isIdeas(text string) bool {
	re := regexp.MustCompile("(?i)#Идеи_покупок|#ИдеиПокупок|#ИнвестИдея")
	return re.FindStringIndex(text) != nil
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
