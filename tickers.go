package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

func init() {
	go watch(func() {
		tmp, err := load()
		if err != nil {
			log.Printf("Can't initialise config: %s", err)
			return
		}
		// tickersMu.Lock()
		// defer tickersMu.Unlock()
		log.Printf("%#v", tmp)
		tickers = tmp
	})
}

const filename = "tickers.json"

var path = filepath.Join(".", filename)

func load() ([]Ticker, error) {
	var (
		result   []Ticker
		err      error
		file     *os.File
		jsonData []byte
	)

	file, err = os.Open(path)
	if err != nil {
		log.Printf("Failed to open file %s: %s", path, err)
		return nil, err
	}
	defer file.Close()

	jsonData, err = ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file %s: %s", path, err)
		return nil, err
	}

	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Printf("Failed to unmarshal file %s: %s", path, err)
		return nil, err
	}

	return result, err
}

func watch(reload func()) {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify write events.
	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Print(event) // Print the event's info.
				_ = event
				reload()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this path for changes.
	if err := w.Add(path); err != nil {
		log.Fatalln(err)
	}

	reload()

	// Start the watching process - it'll check for changes every 1000ms.
	if err := w.Start(1000 * time.Millisecond); err != nil {
		log.Fatalln(err)
	}
}

// Ticker of stock market
type Ticker struct {
	Symbol string `json:"symbol"`
	Title  string `json:"title"`
}

// from https://stockanalysis.com/stocks/
var tickers = []Ticker{}

// func Filter(vs []string, f func(string) bool) []string {
// 	vsf := make([]string, 0)
// 	for _, v := range vs {
// 			if f(v) {
// 					vsf = append(vsf, v)
// 			}
// 	}
// 	return vsf
// }

// func(v string) bool {
// 	return strings.Contains(v, "e")
// }

// GetTickers function
func GetTickers(search string) []Ticker {
	result := make([]Ticker, 0)
	if len(search) > 0 {
		search = strings.ToUpper(search)
		for _, ticker := range tickers {
			if strings.HasPrefix(strings.ToUpper(ticker.Symbol), search) {
				result = append(result, ticker)
				if len(search) == 1 {
					break
				}
			}
		}
	}
	return result
}

// GetExactTicker function
func GetExactTicker(search string) *Ticker {
	var result *Ticker
	if len(search) > 0 {
		search = strings.ToUpper(search)
		for _, ticker := range tickers {
			if strings.ToUpper(ticker.Symbol) == search {
				result = &ticker
				break
			}
		}
	}
	return result
}
