package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/IvanMenshykov/MoonPhase"
	ss "github.com/comerc/segezha4/screenshot"
	"github.com/comerc/segezha4/utils"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: держать запросы от пользователей в очереди, пока выполняется runBackgroundTask

// TODO: источник по ТА https://finviz.com/screener.ashx?v=210&s=ta_p_tlresistance
// TODO: источник по ТА https://ru.investing.com/equities/facebook-inc-technical
// TODO: Тикеры с точкой BRK.B RDS.A
// TODO: не убивать инстанс chrome
// TODO: сохранять id узеров бота для рассылки когда /start
// TODO: badger для tickers
// TODO: подсказки, если неправильные команды в приватном чате
// TODO: параллельная обработка https://gobyexample.ru/worker-pools.html
// TODO: выводить сообщение о лимите по пересылке

// TODO: оптимизация chromedp
// Q: Chrome exits as soon as my Go program finishes
// A: On Linux, chromedp is configured to avoid leaking resources by force-killing any started Chrome child processes. If you need to launch a long-running Chrome instance, manually start Chrome and connect using RemoteAllocator. https://github.com/chromedp/chromedp/blob/dac8c91f6982c771775a2cc1858b1dcc6bb987a3/allocate_test.go

// https://github.com/chromedp/chromedp/issues/297#issuecomment-487833337
// https://github.com/GoogleChrome/chrome-launcher/blob/master/docs/chrome-flags-for-tools.md
// https://devmarkpro.com/chromedp-get-started
// https://github.com/chromedp/chromedp/issues/687
// https://github.com/chromedp/docker-headless-shell/blob/master/README.md

// TODO: упаковать в Docker chromedp https://hub.docker.com/r/chromedp/headless-shell/

// TODO: пересылать ответы для "Andrew Ka2" к "Andrew Ka"
// TODO: автоматизировать пересылку и разделить отчеты "Инвестиции USA Markets"
// TODO: /info tipranks.com LIFE
// TODO: бумажка пробила 9EMA на дневке?
// TODO: запретить повторы за один день для !! !
// TODO: виджет из википедии по названию компании
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
// TODO: добавить биток GBTC
// TODO: выборка с графиками https://finviz.com/screener.ashx?v=212&t=ZM,BA,MU,MS,GE,AA
// TODO: https://stockcharts.com/h-sc/ui?s=$CPCE https://school.stockcharts.com/doku.php?id=market_indicators:put_call_ratio

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	utils.InitTimeoutFactor()

	var (
		// port      = os.Getenv("PORT")
		// publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		chatID = os.Getenv("SEGEZHA4_CHAT_ID") // you must add it to your config vars
		token  = os.Getenv("SEGEZHA4_SECRET")  // you must add it to your config vars
	)
	// webhook := &tb.Webhook{
	// 	Listen:   ":" + port,
	// 	Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	// }
	pref := tb.Settings{
		// URL:    "https://api.bots.mn/telegram/",
		Token: token,
		// Poller: webhook,
		Poller: &tb.LongPoller{Timeout: 10 * time.Minute},
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
					Title:       fmt.Sprintf("%s %s", articleCase.name, ticker.symbol),
					Description: ticker.title,
					HideURL:     true,
					URL:         linkURL,
					ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticker.symbol), // from stockanalysis.com
				}
			} else {
				title := fmt.Sprintf("%s %s", articleCase.name, ticker.symbol)
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
				Text: fmt.Sprintf("/%s %s",
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

*Inline Menu Mode:*
@TickerInfoBot TSLA

*Simple (Batch) Mode:*
#TSLA! #TSLA? #TSLA?? #TSLA?! #TSLA!!
`
			sendText(b, m.Chat.ID, escape(help), false)
		} else if text == "/pause" {
			if isAdmin(m.Sender.ID) {
				pauseDay = time.Now().UTC().Day()
				sendText(b, m.Chat.ID, "set pause", false)
			}
		} else if text == "/reset" {
			if isAdmin(m.Sender.ID) {
				pauseDay = -1
				sendText(b, m.Chat.ID, "set reset", false)
			}
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
		} else if articleCase := hasArticleCase(text); articleCase != nil {
			re := regexp.MustCompile(",|[ ]+")
			payload := re.ReplaceAllString(strings.Trim(m.Payload, " "), " ")
			symbols := strings.Split(payload, " ")
			callbacks := make([]getWhat, 0)
			executed := make([]string, 0)
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
					symbol = symbol[1:]
				}
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				callbacks = append(callbacks, closeWhat(symbol, articleCase))
			}
			sendBatch(b, m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if isEarnings(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)), m.Chat.Type != tb.ChatPrivate)
					continue
				}
				articleCase := GetExactArticleCase("marketwatch")
				result := sendScreenshotForMarketWatch(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendLink(b, m.Chat.ID, articleCase, ticker, m.Chat.Type != tb.ChatPrivate)
				}
			}
		} else if isARKOrWatchList(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			executed = append(executed, "ARK")
			executed = append(executed, "WATCH") // for #Watch_list by @usamarke1
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)), m.Chat.Type != tb.ChatPrivate)
					continue
				}
				articleCase := GetExactArticleCase("finviz")
				result := sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)), m.Chat.Type != tb.ChatPrivate)
					result = true
				}
			}
		} else if isIdeas(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				ticker := GetExactTicker(symbol)
				if ticker == nil {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found`, strings.ToUpper(symbol)), m.Chat.Type != tb.ChatPrivate)
					continue
				}
				articleCase := GetExactArticleCase("finviz")
				result := sendScreenshotForFinviz(b, m.Chat.ID, articleCase, ticker)
				if !result {
					sendText(b, m.Chat.ID, fmt.Sprintf(`\#%s not found on finviz\.com`, strings.ToUpper(symbol)), m.Chat.Type != tb.ChatPrivate)
					result = true
				}
			}
		} else if symbol := hasDots(text); symbol != "" {
			getWhat := closeWhat(symbol, GetExactArticleCase("chart"))
			send(b, m.Chat.ID, m.Chat.Type != tb.ChatPrivate, getWhat())
		} else {
			// simple command mode
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)(\?!|\?\?|\?|!!|!)`)
			matches := re.FindAllStringSubmatch(text, -1)
			if len(matches) == 0 {
				if m.Chat.Type == tb.ChatPrivate {
					sendText(b, m.Chat.ID, escape("Unknown command, please see /help"), false)
				}
				return
			}
			callbacks := make([]getWhat, 0)
			executed := make([]string, 0)
			for _, match := range matches {
				symbol := match[2]
				mode := match[3]
				if utils.Contains(executed, strings.ToUpper(symbol)+mode) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol)+mode)
				switch mode {
				case "?!":
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("marketwatch")))
				case "??":
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("barchart")))
				case "?":
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("stockscores")))
				case "!!":
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("shortvolume")))
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("stockscores")))
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("finviz")))
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("gurufocus")))
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("marketbeat")))
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("tipranks")))
				case "!":
					callbacks = append(callbacks, closeWhat(symbol, GetExactArticleCase("finviz")))
				}
			}
			sendBatch(b, m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		}
	}
	b.Handle(tb.OnText, messageHandler)
	b.Handle(tb.OnPhoto, messageHandler)
	pauseDay = -1
	go runBackgroundTask(b, int64(utils.ConvertToInt(chatID)))
	b.Start()
}

func escapeURL(s string) string {
	re := regexp.MustCompile("[(|)]")
	return re.ReplaceAllString(s, `\$0`)
}

func escape(s string) string {
	re := regexp.MustCompile(`[.|\-|\_|(|)|#|!]`)
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

// func sendScreenshotForPage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
// 	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
// 	screenshot := ss.MakeScreenshotForPage(linkURL, articleCase.x, articleCase.y, articleCase.width, articleCase.height)
// 	if len(screenshot) == 0 {
// 		return false
// 	}
// 	photo := &tb.Photo{
// 		File: tb.FromReader(bytes.NewReader(screenshot)),
// 		Caption: fmt.Sprintf(
// 			`\#%s by [%s](%s)`,
// 			ticker.symbol,
// 			escape(articleCase.name),
// 			linkURL,
// 			// getUserLink(m.Sender),
// 		),
// 	}
// 	_, err := b.Send(
// 		tb.ChatID(chatID),
// 		photo,
// 		&tb.SendOptions{
// 			ParseMode: tb.ModeMarkdownV2,
// 		},
// 	)
// 	screenshot = nil
// 	photo = nil
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }

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
	screenshot := ss.MakeScreenshotForTipRanks(linkURL)
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

// func sendScreenshotForImage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker) bool {
// 	imageURL := fmt.Sprintf(articleCase.imageURL, ticker.symbol)
// 	linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.symbol))
// 	screenshot := ss.MakeScreenshotForImage(imageURL, articleCase.width, articleCase.height)
// 	if len(screenshot) == 0 {
// 		return false
// 	}
// 	photo := &tb.Photo{
// 		File: tb.FromReader(bytes.NewReader(screenshot)),
// 		Caption: fmt.Sprintf(
// 			`\#%s by [%s](%s)`,
// 			ticker.symbol,
// 			escape(articleCase.name),
// 			linkURL,
// 			// getUserLink(m.Sender),
// 		),
// 	}
// 	_, err := b.Send(
// 		tb.ChatID(chatID),
// 		photo,
// 		&tb.SendOptions{
// 			ParseMode: tb.ModeMarkdownV2,
// 		},
// 	)
// 	screenshot = nil
// 	photo = nil
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }

func sendFinvizImage(b *tb.Bot, chatID int64, symbol string, isSlowMode bool) bool {
	if isSlowMode {
		time.Sleep(4 * time.Second) // your bot will not be able to send more than 20 messages per minute to the same group.
	}
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

func sendImage(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker, isSlowMode bool) bool {
	if isSlowMode {
		time.Sleep(4 * time.Second) // your bot will not be able to send more than 20 messages per minute to the same group.
	}
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

func sendLink(b *tb.Bot, chatID int64, articleCase *ArticleCase, ticker *Ticker, isSlowMode bool) {
	if isSlowMode {
		time.Sleep(4 * time.Second) // your bot will not be able to send more than 20 messages per minute to the same group.
	}
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
		escape(utils.GetHost(linkURL)),
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

func sendText(b *tb.Bot, chatID int64, text string, isSlowMode bool) {
	if isSlowMode {
		time.Sleep(4 * time.Second) // your bot will not be able to send more than 20 messages per minute to the same group.
	}
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

var pauseDay int

func runBackgroundTask(b *tb.Bot, chatID int64) {
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		utc := t.UTC()
		w := utc.Weekday()
		if w == 6 || w == 0 {
			continue
		}
		d := utc.Day()
		if d == pauseDay {
			continue
		} else if pauseDay > -1 {
			pauseDay = -1 // reset
		}
		h := utc.Hour()
		m := utc.Minute()
		s := utc.Second()
		const (
			delta  = 30
			summer = 1
		)
		if h == 14-summer && m >= 30 || h > 14-summer && h < 21-summer || h == 21-summer && m < delta {
			if m%delta == 0 && s == 15 {
				if h == 14-summer && m >= 30 {
					moon := MoonPhase.New(t)
					isFullMoon := int(math.Floor((moon.Phase()+0.0625)*8)) == 4
					if isFullMoon {
						sendText(b, chatID, escape("🌕 #FullMoon"), false)
					}
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

		// if s%10 == 0 {
		// 	go func(t time.Time) {
		// 		chatID2 := -1001374011821 // ticker_info_test_channel_1
		// 		// chatID2 := -1001211314640 // ticker_info_test_channel_2
		// 		msg, err1 := b.Send(
		// 			tb.ChatID(chatID2),
		// 			"send text "+t.String(),
		// 		)
		// 		if err1 != nil {
		// 			log.Println(err1)
		// 		}
		// 		time.Sleep(5 * time.Second)
		// 		_, err2 := b.Edit(
		// 			msg,
		// 			"*edit text* "+escape(fmt.Sprintf(`https://t.me/%s/%d`, msg.Chat.Username, msg.ID)),
		// 			tb.ModeMarkdownV2,
		// 		)
		// 		if err2 != nil {
		// 			log.Println(err2)
		// 		}
		// 	}(t)
		// }
	}
}

func sendBarChart(b *tb.Bot, chatID int64, symbol string) bool {
	volume, height, tag := func() (string, string, string) {
		if strings.HasPrefix(symbol, "$") {
			return "0", "O", ""
		}
		return "total", "X", "#"
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

func isEarnings(text string) bool {
	re := regexp.MustCompile("#ОТЧЕТ")
	return re.FindStringIndex(text) != nil
}

func isARKOrWatchList(text string) bool {
	re := regexp.MustCompile("#ARK|#Watch_list")
	return re.FindStringIndex(text) != nil
}

func isIdeas(text string) bool {
	re := regexp.MustCompile("(?i)#Идеи_покупок|#ИдеиПокупок|#ИнвестИдея")
	return re.FindStringIndex(text) != nil
}

func hasDots(text string) string {
	re := regexp.MustCompile(`(\x{1F7E2}\x{1F7E2}|\x{1F534}\x{1F534}) ([A-Za-z]+)`) // green / red dots
	matches := re.FindAllStringSubmatch(text, -1)
	if len(matches) == 1 {
		return matches[0][2]
	}
	return ""
}

func isAdmin(ID int) bool {
	s := os.Getenv("SEGEZHA4_ADMIN_USER_IDS")
	ids := strings.Split(s, ",")
	return utils.Contains(ids, fmt.Sprintf("%d", ID))
}

type getWhat func() interface{}

func closeWhat(symbol string, articleCase *ArticleCase) getWhat {
	return func() interface{} {
		tag := func() string {
			if strings.HasPrefix(symbol, "$") {
				return ""
			}
			return "#"
		}()
		// TODO: пополнять базу тикеров и индексов для inline mode
		if tag == "#" {
			if isNotFoundTicker(symbol) {
				return fmt.Sprintf("%s not found", strings.ToUpper(tag+symbol))
			}
		}
		// TODO: not found for $symbol
		var result interface{}
		linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(symbol))
		defer utils.Elapsed(linkURL)()
		switch articleCase.screenshotMode {
		case ScreenshotModeImage:
			imageURL := fmt.Sprintf(articleCase.imageURL, strings.ToLower(symbol), time.Now().Unix())
			result = &tb.Photo{
				File:    tb.FromURL(imageURL),
				Caption: getCaption(linkURL, tag+symbol),
			}
		case ScreenshotModeFinviz:
			screenshot := ss.MakeScreenshotForFinviz(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeMarketWatch:
			screenshot := ss.MakeScreenshotForMarketWatch(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeCathiesArk:
			screenshot := ss.MakeScreenshotForCathiesArk(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeGuruFocus:
			screenshot := ss.MakeScreenshotForGuruFocus(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeMarketBeat:
			screenshot := ss.MakeScreenshotForMarketBeat(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeTipRanks:
			screenshot := ss.MakeScreenshotForTipRanks(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(linkURL, tag+symbol),
				}
			}
		case ScreenshotModeBarChart:
			volume, height := func() (string, string) {
				if strings.HasPrefix(symbol, "$") {
					return "0", "O"
				}
				return "total", "X"
			}()
			linkURL := "https://www.barchart.com/stocks/quotes/%s/technical-chart%s?plot=CANDLE&volume=%s&data=I:15&density=%[4]s&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(50);EXPMA(20);EXPMA(200);WMA(9);EXPMA(500)&sym=%[1]s&grid=1&height=500&studyheight=200"
			screenshot := ss.MakeScreenshotForBarChart(fmt.Sprintf(linkURL, symbol, "/fullscreen", volume, height))
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(fmt.Sprintf(linkURL, symbol, "", volume, height), tag+symbol),
				}
			}
		}
		if result == nil {
			result = getEmptyText(linkURL, articleCase.description, tag+symbol)
		}
		return result
	}
}

func getCaption(linkURL, tagSymbol string) string {
	return fmt.Sprintf("%s by [%s](%s)",
		escape(strings.ToUpper(tagSymbol)),
		escape(utils.GetHost(linkURL)),
		escapeURL(linkURL),
	)
}

func getEmptyText(linkURL, description, tagSymbol string) string {
	return fmt.Sprintf("%s %s[%s](%s)",
		escape(strings.ToUpper(tagSymbol)),
		escape(by(description)),
		escape(utils.GetHost(linkURL)),
		escapeURL(linkURL),
	)
}

// **** параллельная обработка

type ParallelResult struct {
	what       interface{}
	isReceived bool
	isSent     bool
}

func sendBatch(b *tb.Bot, chatID int64, isPrivateChat bool, callbacks []getWhat) {
	defer utils.Elapsed("sendBatch")()
	if len(callbacks) == 0 {
		return
	}
	done := make(chan bool)
	results := make([]ParallelResult, len(callbacks))
	threads := utils.ConvertToInt(os.Getenv("SEGEZHA4_THREADS"))
	if threads == 0 {
		threads = 1
	}
	var tokens = make(chan struct{}, threads) // ограничение количества горутин
	var mu sync.Mutex
	receivedCount := 0
	for i, cb := range callbacks {
		tokens <- struct{}{} // захват маркера
		go func(i int, cb getWhat) {
			what := cb()
			<-tokens // освобождение маркера
			{
				mu.Lock()
				defer mu.Unlock()
				results[i] = ParallelResult{
					what:       what,
					isReceived: true,
				}
				receivedCount = receivedCount + 1
				if receivedCount == len(callbacks) {
					for i, r := range results {
						func(i int, r ParallelResult) {
							if !r.isSent {
								send(b, chatID, isPrivateChat, r.what)
								results[i].isSent = true
							}
						}(i, r)
					}
					done <- true
				} else {
					isAllPreviosReceived := true
					for _, r := range results[:i] {
						if !r.isReceived {
							isAllPreviosReceived = false
							break
						}
					}
					if isAllPreviosReceived {
						for i, r := range results[:i+1] {
							func(i int, r ParallelResult) {
								if !r.isSent {
									send(b, chatID, isPrivateChat, r.what)
									results[i].isSent = true
								}
							}(i, r)
						}
					}
				}
			}
		}(i, cb)
	}
	<-done
}

var lastSendByGroup = make(map[int64]time.Time)

const pause = 4 * time.Second

func send(b *tb.Bot, chatID int64, isPrivateChat bool, what interface{}) {
	if !isPrivateChat {
		// your bot will not be able to send more than 20 messages per minute to the same group.
		lastSend := lastSendByGroup[chatID]
		diff := time.Since(lastSend)
		if diff < pause {
			time.Sleep(pause)
		}
		lastSendByGroup[chatID] = time.Now()
	}
	_, err := b.Send(
		tb.ChatID(chatID),
		what,
		&tb.SendOptions{
			ParseMode:             tb.ModeMarkdownV2,
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func isNotFoundTicker(symbol string) bool {
	// TODO: реализация пополнения тикеров
	ticker := GetExactTicker(symbol)
	return ticker == nil
}

func hasArticleCase(text string) *ArticleCase {
	if text != "" {
		text = strings.ToUpper(text)
		for _, articleCase := range ArticleCases {
			command := fmt.Sprintf("/%s ", strings.ToUpper(articleCase.name))
			if strings.HasPrefix(text, command) {
				return &articleCase
			}
		}
	}
	return nil
}
