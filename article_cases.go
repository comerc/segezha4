package main

import (
	"strings"
)

// ArticleCase struct
type ArticleCase struct {
	name     string
	linkURL  string
	imageURL string
	hasGift  bool
}

// TODO: чаще используемые перемещать наверх
// TODO: конфигурировать через интерфейс

// ArticleCases slice
var ArticleCases = []ArticleCase{
	{name: "tradingview.com", linkURL: "https://ru.tradingview.com/symbols/%s"},
	{name: "finviz.com", linkURL: "https://finviz.com/quote.ashx?t=%s", hasGift: true},
	{name: "cathiesark.com", linkURL: "https://cathiesark.com/ark-combined-holdings-of-%s"},
	{name: "marketwatch.com", linkURL: "https://marketwatch.com/investing/stock/%s"},
	{name: "stockscores.com", linkURL: "https://stockscores.com/charts/charts/?ticker=%s",
		hasGift:  true,
		imageURL: "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=180&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1800&ChartHeight=300&LogScale=None&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=None&Indicator3=MACD&Indicator4=AccDist&endDate=2021-2-17&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// imageURL: "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=180&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=920&ChartHeight=300&LogScale=None&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=None&Indicator3=MACD&Indicator4=AccDist&endDate=2021-2-17&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// imageURL: "https://stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=180&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1600&ChartHeight=480&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=MACD&Indicator3=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// imageURL: "https://stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=100&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=860&ChartHeight=480&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=MACD&Indicator3=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// linkURL: "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=100&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=860&ChartHeight=480&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=50&Indicator1=Momentum&Indicator2=RSI&Indicator3=MACD&Indicator4=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// linkURL: "https://stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=365&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1920&ChartHeight=430&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=50&Indicator1=Momentum&Indicator2=RSI&Indicator3=MACD&Indicator4=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
	},
	{name: "shortvolume.com", linkURL: "https://shortvolume.com/?t=%s",
		hasGift:  true,
		imageURL: "https://shortvolume.com/chart_engine/draw_chart.php?Symbol=%s&TimeRange=100"},
	{name: "marketbeat.com", linkURL: "https://marketbeat.com/stocks/%s"},
	{name: "earningswhispers.com", linkURL: "https://earningswhispers.com/stocks/%s"},
	// {name: "tipranks.com", linkURL: "https://tipranks.com/stocks/%s/forecast"},
	{name: "barchart.com", linkURL: "https://barchart.com/stocks/quotes/%s/overview"},
	{name: "gurufocus.com", linkURL: "https://gurufocus.com/stock/%s/summary"},
	{name: "stockrow.com", linkURL: "https://stockrow.com/%s"},
	{name: "stockanalysis.com", linkURL: "https://stockanalysis.com/stocks/%s/"},
	{name: "finasquare.com", linkURL: "https://www.finasquare.com/stocks/%s/company-info/overview"},
}

// GetExactArticleCase function
func GetExactArticleCase(search string) *ArticleCase {
	var result *ArticleCase
	if len(search) > 0 {
		search = strings.ToUpper(search)
		for _, articleCase := range ArticleCases {
			if strings.ToUpper(articleCase.name) == search {
				result = &articleCase
				break
			}
		}
	}
	return result
}
