package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/IvanMenshykov/MoonPhase"
	ss "github.com/comerc/segezha4/screenshot"
	"github.com/comerc/segezha4/utils"
	"github.com/dgraph-io/badger"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO: #AMD!? выдавать "Unknown command, please see /help"

// TODO: badger для tickers и добавлять, когда "not found"

// TODO: бумажка пробила 9EMA на дневке?

// TODO: /intro

// TODO: короткие команды /mw /fv /mb /ew

// TODO: "Самые обсуждаемые акции на форумах" - выдавать графики

// TODO: https://stockcharts.com/h-sc/ui?s=$CPCE https://school.stockcharts.com/doku.php?id=market_indicators:put_call_ratio

// TODO: запретить команды для публичных чатов

// TODO: /crypto dogeusd btcusd ethusd xrpusd bchusd ltcusd xmrusd (https://www.marketwatch.com/investing/cryptocurrency/btcusd)

// TODO: в @teslaholics2 при клике по ссылке внутри сообщения /help - /help@TickerInfoBot
// TODO: держать запросы от пользователей в очереди, пока выполняется runBackgroundTask

// TODO: источник по ТА https://finviz.com/screener.ashx?v=210&s=ta_p_tlresistance
// TODO: источник по ТА https://ru.investing.com/equities/facebook-inc-technical
// TODO: Тикеры с точкой BRK.B RDS.A (finviz заменяет на "-")
// TODO: подсказки, если неправильные команды в приватном чате
// TODO: демо всех тикеров в приватном чате
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
// TODO: подключить tradingview.com
// TODO: добавить тайм-фрейм #BABA?15M
// TODO: добавить медленную скользящую #BABA?50EMA / 100EMA / 200EMA
// TODO: выборка с графиками https://finviz.com/screener.ashx?v=212&t=ZM,BA,MU,MS,GE,AA

var (
	db *badger.DB
	b  *tb.Bot
)

const help = `*Commands:*
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

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	utils.InitTimeoutFactor()

	{
		path := filepath.Join(".", "data")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
		var err error
		db, err = badger.Open(badger.DefaultOptions(path))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer db.Close()

	var (
		// port      = os.Getenv("PORT")
		// publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		chatID  = os.Getenv("SEGEZHA4_CHAT_ID") // you must add it to your config vars
		token   = os.Getenv("SEGEZHA4_SECRET")  // you must add it to your config vars
		pingURL = os.Getenv("SEGEZHA4_PING_URL")
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
	{
		var err error
		b, err = tb.NewBot(pref)
		if err != nil {
			log.Fatal(err)
		}
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
			linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(ticker.Symbol))
			var result *tb.ArticleResult
			if i == 0 {
				result = &tb.ArticleResult{
					Title:       fmt.Sprintf("%s %s", articleCase.name, ticker.Symbol),
					Description: ticker.Title,
					HideURL:     true,
					URL:         linkURL,
					ThumbURL:    fmt.Sprintf("https://storage.googleapis.com/iexcloud-hl37opg/api/logos/%s.png", ticker.Symbol), // from stockanalysis.com
				}
			} else {
				title := fmt.Sprintf("%s %s", articleCase.name, ticker.Symbol)
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
					ticker.Symbol,
				),
				DisablePreview: true,
			})
			result.SetResultID(ticker.Symbol + "=" + articleCase.name)
			results[i] = result
		}
		if err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60,
		}); err != nil {
			log.Println(err)
		}
	})
	messageHandler := func(m *tb.Message) {
		log.Println("****")
		log.Println("LastEdit:", m.LastEdit)
		if m.Sender != nil {
			log.Println("Username:", m.Sender.Username)
			log.Println("FirstName:", m.Sender.FirstName)
			log.Println("LastName:", m.Sender.LastName)
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
				send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhatMarketWatchIDs(tab))
				return
			}
		}
		if text == "/start" || text == "/help" {
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, escape(help))
		} else if text == "/stat" && isAdmin(m.Sender.ID) {
			var totalKeys, totalValues int64
			if err := db.View(func(txn *badger.Txn) error {
				opts := badger.DefaultIteratorOptions
				opts.PrefetchSize = 10
				it := txn.NewIterator(opts)
				defer it.Close()
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					k := item.Key()
					totalKeys += 1
					if err := item.Value(func(v []byte) error {
						key := int64(bytesToUint64(k))
						val := int64(bytesToUint64(v))
						totalValues += val
						log.Print(key, val)
						return nil
					}); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				log.Print(err)
			}
			log.Printf("keys: %d values: %d", totalKeys, totalValues)
		} else if text == "/pause" && isAdmin(m.Sender.ID) {
			pauseDay = time.Now().UTC().Day()
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, "pause")
		} else if text == "/reset" && isAdmin(m.Sender.ID) {
			pauseDay = -1
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, "reset")
		} else if text == "/bb" {
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhatFinvizBB())
		} else if text == "/vix" {
			getWhat := closeWhat("$VIX", GetExactArticleCase("barchart"))
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhat())
		} else if text == "/spy" {
			getWhat := closeWhat("SPY", GetExactArticleCase("barchart"))
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhat())
		} else if text == "/index" {
			callbacks := make([]getWhat, 0)
			articleCase := GetExactArticleCase("barchart")
			callbacks = append(callbacks, closeWhat("$INX", articleCase))
			callbacks = append(callbacks, closeWhat("$NASX", articleCase))
			callbacks = append(callbacks, closeWhat("$DOWI", articleCase))
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if text == "/volume" {
			callbacks := make([]getWhat, 0)
			articleCase := GetExactArticleCase("barchart")
			callbacks = append(callbacks, closeWhat("SPY", articleCase))
			callbacks = append(callbacks, closeWhat("QQQ", articleCase))
			callbacks = append(callbacks, closeWhat("DOW", articleCase))
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if text == "/map" {
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhatFinvizMap())
		} else if text == "/fear" {
			send(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, getWhatFear())
		} else if articleCase := hasArticleCase(text); articleCase != nil {
			re := regexp.MustCompile(",|[ ]+")
			payload := re.ReplaceAllString(strings.Trim(m.Payload, " "), " ")
			symbols := strings.Split(payload, " ")
			executed := make([]string, 0)
			callbacks := make([]getWhat, 0)
			for _, symbol := range symbols {
				if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") && !isBarChart(text) {
					symbol = symbol[1:]
				}
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				callbacks = append(callbacks, closeWhat(symbol, articleCase))
			}
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if isEarnings(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			articleCase := GetExactArticleCase("marketwatch")
			callbacks := make([]getWhat, 0)
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				callbacks = append(callbacks, closeWhat(symbol, articleCase))
			}
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if isARKOrWatchList(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			executed = append(executed, "ARK")
			if m.Chat.Username == "usamarke1" {
				executed = append(executed, "WATCH") // for #Watch_list
			}
			articleCase := GetExactArticleCase("finviz")
			callbacks := make([]getWhat, 0)
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				callbacks = append(callbacks, closeWhat(symbol, articleCase))
			}
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if isIdeas(text) {
			re := regexp.MustCompile(`(^|[^A-Za-z])\$([A-Za-z]+)`)
			matches := re.FindAllStringSubmatch(text, -1)
			executed := make([]string, 0)
			articleCase := GetExactArticleCase("finviz")
			callbacks := make([]getWhat, 0)
			for _, match := range matches {
				symbol := match[2]
				if utils.Contains(executed, strings.ToUpper(symbol)) {
					continue
				}
				executed = append(executed, strings.ToUpper(symbol))
				callbacks = append(callbacks, closeWhat(symbol, articleCase))
			}
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		} else if symbol := hasDots(text); symbol != "" {
			getWhat := closeWhat(symbol, GetExactArticleCase("chart"))
			send(m.Chat.ID, m.Chat.Type != tb.ChatPrivate, getWhat())
		} else {
			// simple command mode
			re := regexp.MustCompile(`(^|[^A-Za-z])#([A-Za-z]+)(\?!|\?\?|\?|!!|!)`)
			matches := re.FindAllStringSubmatch(text, -1)
			if len(matches) == 0 {
				if m.Chat.Type == tb.ChatPrivate {
					send(m.Chat.ID, m.Chat.Type != tb.ChatPrivate, escape("Unknown command, please see /help"))
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
			sendBatch(m.Chat.ID, m.Chat.Type == tb.ChatPrivate, callbacks)
		}
	}
	b.Handle(tb.OnText, messageHandler)
	b.Handle(tb.OnPhoto, messageHandler)
	pauseDay = -1
	go runBackgroundTask(b, int64(utils.ConvertToInt(chatID)), pingURL)
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

// TODO: replace escape() to escapeMarkdown()
// func escapeMarkdown(s string) string {
// 	// You can escape the following characters:
// 	// Asterisk \*
// 	// Underscore \_
// 	// Curly braces \{ \}
// 	// Square brackets \[ \]
// 	// Brackets \( \)
// 	// Hash \#
// 	// Plus \+
// 	// Minus \-
// 	// Period \.
// 	// Exclamation point \!
// 	a := []string{"*", `\_`, "{", "}", `\[`, `\]`, `\(`, `\)`, "#", "+", `\-`, ".", "!"}
// 	re := regexp.MustCompile("[" + strings.Join(a, "|") + "]")
// 	return re.ReplaceAllString(s, `\$0`)
// }

// func getUserLink(u *tb.User) string {
// 	if u.Username != "" {
// 		return fmt.Sprintf("@%s", u.Username)
// 	}
// 	return fmt.Sprintf("[%s](tg://user?id=%d)", u.FirstName, u.ID)
// }

func by(s string) string {
	if s == "" {
		return "by "
	}
	return s + " by "
}

var (
	pauseDay   int
	currentDay int
)

func runBackgroundTask(b *tb.Bot, chatID int64, pingURL string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for t := range ticker.C {
		utc := t.UTC()
		s := utc.Second()
		if s == 0 {
			go func() {
				netClient := &http.Client{
					Timeout: 10 * time.Second,
				}
				isAlarm := false
				response, err := netClient.Get(fmt.Sprintf("%s?rand=%d", pingURL, time.Now().Unix()))
				if err != nil {
					log.Printf("netClient.Get(pingURL): %s", err)
					isAlarm = true
				} else if response.StatusCode != 200 {
					log.Print("netClient.Get(pingURL): response.StatusCode != 200")
					isAlarm = true
				}
				if isAlarm {
					s := os.Getenv("SEGEZHA4_ADMIN_USER_IDS")
					IDs := strings.Split(s, ",")
					for _, ID := range IDs {
						_, err := b.Send(
							tb.ChatID(utils.ConvertToInt(ID)),
							fmt.Sprintf("Not responsed %s", pingURL),
						)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}()
		}
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
		if d != currentDay {
			currentDay = d
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
		h := utc.Hour()
		m := utc.Minute()
		const (
			delta  = 30
			summer = 1
		)
		callbacks := make([]getWhat, 0)
		if h == 14-summer && m >= 30 || h > 14-summer && h < 21-summer || h == 21-summer && m < delta {
			if m%delta == 0 && s == 15 {
				if h == 14-summer && m >= 30 {
					moon := MoonPhase.New(t)
					isFullMoon := int(math.Floor((moon.Phase()+0.0625)*8)) == 4
					if isFullMoon {
						callbacks = append(callbacks, getWhatFullMoon)
					}
					callbacks = append(callbacks, getWhatFear)
				}
				if h >= 15-summer {
					callbacks = append(callbacks, getWhatFinvizBB)
					callbacks = append(callbacks, getWhatFinvizMap)
				}
				callbacks = append(callbacks, closeWhat("$VIX", GetExactArticleCase("barchart")))
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabUS))
				if h >= 8 && h <= 17 {
					callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabEurope))
				}
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabRates))
			}
		} else if m == 0 && s == 15 {
			if h >= 8 && h <= 17 {
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabEurope))
			}
			// SPB работает с 7 утра (MSK)
			if h >= 4 && h <= 9 {
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabAsia))
			}
			if h >= 4 && h <= 14-summer {
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabRates))
				callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabFutures))
			}
			// callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabFX))
			// callbacks = append(callbacks, closeWhatMarketWatchIDs(ss.MarketWatchTabCrypto))
		}
		sendBatch(chatID, false, callbacks)

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

func getWhatFinvizMap() interface{} {
	linkURL := "https://finviz.com/map.ashx?t=sec"
	defer utils.Elapsed(linkURL)()
	caption := getCaption("#map", "", linkURL)
	screenshot := ss.MakeScreenshotForFinvizMap(linkURL)
	if len(screenshot) == 0 {
		return caption
	}
	return &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(screenshot)),
		Caption: caption,
	}
}

func getWhatFullMoon() interface{} {
	return escape("🌕 #FullMoon")
}

func getWhatFear() interface{} {
	linkURL := "https://money.cnn.com/data/fear-and-greed/"
	defer utils.Elapsed(linkURL)()
	caption := getCaption("#fear", "", linkURL)
	screenshot := ss.MakeScreenshotForFear(linkURL)
	if len(screenshot) == 0 {
		return caption
	}
	return &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(screenshot)),
		Caption: caption,
	}
}

func getWhatFinvizBB() interface{} {
	linkURL := "https://finviz.com/"
	defer utils.Elapsed(linkURL)()
	caption := getCaption("#bb", "Bull or Bear", linkURL)
	screenshot := ss.MakeScreenshotForFinvizBB(linkURL)
	if len(screenshot) == 0 {
		return caption
	}
	return &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(screenshot)),
		Caption: caption,
	}
}

func getWhatMarketWatchIDs(tab ss.MarketWatchTab) interface{} {
	linkURL := "https://www.marketwatch.com/"
	defer utils.Elapsed(linkURL + tab)()
	caption := getCaption("#"+tab, "", linkURL)
	screenshot := ss.MakeScreenshotForMarketWatchIDs(linkURL, tab)
	if len(screenshot) == 0 {
		return caption
	}
	return &tb.Photo{
		File:    tb.FromReader(bytes.NewReader(screenshot)),
		Caption: caption,
	}
}

func isEarnings(text string) bool {
	re := regexp.MustCompile("#ОТЧЕТ") // TODO: #отчетность by @MarketTwits
	return re.FindStringIndex(text) != nil
}

func isARKOrWatchList(text string) bool {
	re := regexp.MustCompile("#ARK Trading Desk|#Watch_list")
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
	IDs := strings.Split(s, ",")
	return utils.Contains(IDs, fmt.Sprintf("%d", ID))
}

type getWhat func() interface{}

func closeWhat(symbol string, articleCase *ArticleCase) getWhat {
	return func() interface{} {
		tag := func() string {
			if strings.HasPrefix(symbol, "$") { // для isBarChart
				return ""
			}
			return "#"
		}()
		// TODO: пополнять базу тикеров и индексов для inline mode
		var ticker *Ticker
		if tag == "#" {
			ticker = GetExactTicker(symbol)
			if ticker == nil {
				return fmt.Sprintf("%s not found", escape(strings.ToUpper(tag+symbol)))
			}
		} else {
			// TODO: not found for $symbol
		}
		var result interface{}
		linkURL := fmt.Sprintf(articleCase.linkURL, strings.ToLower(symbol))
		defer utils.Elapsed(linkURL)()
		switch articleCase.screenshotMode {
		case ScreenshotModeImage:
			imageURL := fmt.Sprintf(articleCase.imageURL, strings.ToLower(symbol), time.Now().Unix())
			result = &tb.Photo{
				File:    tb.FromURL(imageURL),
				Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
			}
		case ScreenshotModeFinviz:
			screenshot := ss.MakeScreenshotForFinviz(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeMarketWatch:
			screenshot := ss.MakeScreenshotForMarketWatch(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeCathiesArk:
			screenshot := ss.MakeScreenshotForCathiesArk(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeGuruFocus:
			screenshot := ss.MakeScreenshotForGuruFocus(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeMarketBeat:
			screenshot := ss.MakeScreenshotForMarketBeat(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeTipRanks:
			screenshot := ss.MakeScreenshotForTipRanks(linkURL)
			if len(screenshot) != 0 {
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		case ScreenshotModeBarChart:
			volume, height := func() (string, string) {
				if strings.HasPrefix(symbol, "$") {
					return "0", "O"
				}
				return "total", "X"
			}()
			srcURL := "https://www.barchart.com/stocks/quotes/%s/technical-chart%s?plot=CANDLE&volume=%s&data=I:15&density=%[4]s&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(50);EXPMA(20);EXPMA(200);WMA(9);EXPMA(500)&sym=%[1]s&grid=1&height=500&studyheight=200"
			dscURL := fmt.Sprintf(srcURL, symbol, "/fullscreen", volume, height)
			screenshot := ss.MakeScreenshotForBarChart(dscURL)
			if len(screenshot) != 0 {
				linkURL := fmt.Sprintf(srcURL, symbol, "", volume, height)
				result = &tb.Photo{
					File:    tb.FromReader(bytes.NewReader(screenshot)),
					Caption: getCaption(strings.ToUpper(tag+symbol), "", linkURL),
				}
			}
		}
		if result == nil {
			description := func() string {
				if articleCase.name == ArticleCases[0].name && ticker != nil {
					return ticker.Title
				}
				return articleCase.description
			}()
			result = getCaption(strings.ToUpper(tag+symbol), description, linkURL)
		}
		return result
	}
}

func getCaption(tagSymbol, description, linkURL string) string {
	return fmt.Sprintf("%s %s[%s](%s)",
		escape(tagSymbol),
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

func sendBatch(chatID int64, isPrivateChat bool, callbacks []getWhat) {
	if len(callbacks) == 0 {
		return
	}
	defer utils.Elapsed("sendBatch")()
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
								send(chatID, isPrivateChat, r.what)
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
									send(chatID, isPrivateChat, r.what)
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

const pause = 3 * time.Second

func send(chatID int64, isPrivateChat bool, what interface{}) {
	if isPrivateChat {
		increment(chatID)
	} else {
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

func closeWhatMarketWatchIDs(tab ss.MarketWatchTab) getWhat {
	return func() interface{} { return getWhatMarketWatchIDs(tab) }
}

func isBarChart(text string) bool {
	return strings.HasPrefix(strings.ToUpper(text), "/BARCHART ")
}

// **** db routines

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Merge function to add two uint64 numbers
func add(existing, new []byte) []byte {
	return uint64ToBytes(bytesToUint64(existing) + bytesToUint64(new))
}

func increment(chatID int64) {
	key := uint64ToBytes(uint64(chatID))
	m := db.GetMergeOperator(key, add, 200*time.Millisecond)
	defer m.Stop()
	err := m.Add(uint64ToBytes(1))
	if err != nil {
		log.Printf("increment() chatID: %d %s", chatID, err)
	}
}
