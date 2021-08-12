package import_tickers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Ticker struct {
	Name         interface{} `json:"name"`
	TickerSymbol interface{} `json:"ticker_symbol"`
	CanonicalUrl string      `json:"canonical_url"`
}

type Meta struct {
	TotalRecords int `json:"total_records"`
}

type Src struct {
	Data []*Ticker `json:"data"`
	Meta Meta      `json:"meta"`
}

var allTickers []*Ticker

const limit = 44

type Dst struct {
	Symbol       string
	Title        string
	SimplyWallSt string
}

func Run() bool {
	allTickers = make([]*Ticker, 0)
	steps := 0
	for step := 0; step <= steps; step++ {
		time.Sleep(time.Duration(getRand(1, 4)) * time.Second)
		log.Print("step ", step)
		offset := step * limit
		totalRecords := getData(offset)
		steps = totalRecords / limit
	}
	if steps < 200 {
		log.Print("error: steps < 200")
		return false
	}
	if len(allTickers) == 0 {
		log.Print("error: allTickers is empty")
		return false
	}
	result := make([]*Dst, 0)
	for _, ticker := range allTickers {
		result = append(result, &Dst{
			Symbol:       fmt.Sprintf("%v", ticker.TickerSymbol),
			Title:        fmt.Sprintf("%v", ticker.Name),
			SimplyWallSt: ticker.CanonicalUrl,
		})
	}
	file, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Print(err)
		return false
	}
	err = ioutil.WriteFile("tickers.json", file, 0644)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func getData(offset int) int {
	requestBody, err := json.Marshal(map[string]interface{}{
		"id":                 0,
		"no_result_if_limit": false,
		"offset":             offset,
		"rules":              `[["order_by","name","asc"],["primary_flag","=",true],["grid_visible_flag","=",true],["market_cap","is_not_null"],["is_fund","=",false],["aor",[["country_name","in",["us"]]]]]`,
		"size":               limit,
		"state":              "read",
	})
	if err != nil {
		log.Print(err)
		return 0
	}
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	request, err := http.NewRequest("POST", "https://api.simplywall.st/api/grid/filter?include=grid,score", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Print(err)
		return 0
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("accept-language", "en")
	request.Header.Add("cache-control", "no-cache")
	request.Header.Add("content-type", "application/json")
	request.Header.Add("origin", "https://simplywall.st")
	request.Header.Add("pragma", "no-cache")
	request.Header.Add("referer", "https://simplywall.st")
	// request.Header.Add("sec-fetch-dest", "empty")
	// request.Header.Add("sec-fetch-mode", "cors")
	// request.Header.Add("sec-fetch-site", "same-site")
	request.Header.Add("user-agent", "Mozilla/5.0")
	// request.Header.Add("x-requested-with", "sws-services/mono-v1.14.33")
	response, err := client.Do(request)
	// response, err := http.Post(
	// 	"https://api.simplywall.st/api/grid/filter?include=grid,score",
	// 	"application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Print(err)
		return 0
	}
	defer response.Body.Close()
	var src Src
	if err := json.NewDecoder(response.Body).Decode(&src); err != nil {
		log.Print(err)
		return 0
	}
	allTickers = append(allTickers, src.Data...)
	return src.Meta.TotalRecords
}

func getRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
