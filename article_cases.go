package main

// import (
// 	"strings"
// )

// ArticleCase struct
// type ArticleCase struct {
// 	name string
// 	url  string
// }

// TODO: чаще используемые перемещать наверх
// TODO: конфигурировать через интерфейс

// TODO: map https://habr.com/ru/post/457728/
// m := make(map[key_type]value_type)
// m := new(map[key_type]value_type)
// var m map[key_type]value_type
// m := map[key_type]value_type{key1: val1, key2: val2}

// ArticleCases map
var ArticleCases = make(map[string]string)

// InitializeArticleCases func
func init() {
	ArticleCases["finviz.com"] = "https://finviz.com/quote.ashx?t=%s"
	ArticleCases["cathiesark.com"] = "https://cathiesark.com/ark-combined-holdings-of-%s"
	ArticleCases["marketwatch.com"] = "https://marketwatch.com/investing/stock/%s"
	ArticleCases["stockscores.com"] =
		"https://stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=100&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=860&ChartHeight=480&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=50&Indicator1=RSI&Indicator2=MACD&Indicator3=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen"
	ArticleCases["shortvolume.com"] = "https://shortvolume.com/chart_engine/draw_chart.php?Symbol=%s&TimeRange=100"
	ArticleCases["marketbeat.com"] = "https://marketbeat.com/stocks/%s"
	ArticleCases["earningswhispers.com"] = "https://earningswhispers.com/stocks/%s"
	// ArticleCases["tipranks.com"] = "https://tipranks.com/stocks/%s/forecast"
	ArticleCases["barchart.com"] = "https://barchart.com/stocks/quotes/%s/overview"
	ArticleCases["gurufocus.com"] = "https://gurufocus.com/stock/%s/summary"
	ArticleCases["stockrow.com"] = "https://stockrow.com/%s"
	ArticleCases["stockanalysis.com"] = "https://stockanalysis.com/stocks/%s/"
	ArticleCases["finasquare.com"] = "https://www.finasquare.com/stocks/%s/company-info/overview"
}

// GetExactArticleCase function
// func GetExactArticleCase(search string) *ArticleCase {
// 	var result *ArticleCase
// 	if len(search) > 0 {
// 		search = strings.ToUpper(search)
// 		for _, articleCase := range ArticleCases {
// 			if articleCase.name == search {
// 				result = &articleCase
// 				break
// 			}
// 		}
// 	}
// 	return result
// }
