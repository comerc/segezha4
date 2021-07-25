package import_tickers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Ticker struct {
	Name         string `json:"name"`
	CanonicalUrl string `json:"canonical_url"`
	TickerSymbol string `json:"ticker_symbol"`
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

func Run() {
	steps := 0
	for step := 0; step <= steps; step++ {
		offset := step * limit
		totalRecords := getData(offset)
		steps = totalRecords / limit
	}
	result := make([]*Dst, 0)
	for _, ticker := range allTickers {
		result = append(result, &Dst{
			Symbol:       ticker.TickerSymbol,
			Title:        ticker.Name,
			SimplyWallSt: ticker.CanonicalUrl,
		})
	}
	file, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Print(err)
		return
	}
	err = ioutil.WriteFile("tickers.json", file, 0644)
	if err != nil {
		log.Print(err)
		return
	}
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
		log.Fatal(err)
	}
	response, err := http.Post(
		"https://api.simplywall.st/api/grid/filter",
		"application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Print(err)
		return 0
	}
	defer response.Body.Close()
	var src Src
	json.NewDecoder(response.Body).Decode(&src)
	allTickers = append(allTickers, src.Data...)
	return src.Meta.TotalRecords
}
