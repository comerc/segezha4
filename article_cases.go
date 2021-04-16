package main

import (
	"strings"
)

// ScreenshotMode description
type ScreenshotMode string

// ScreenshotMode variants
const (
	// ScreenshotModePage        ScreenshotMode = "Page"
	ScreenshotModeImage       ScreenshotMode = "Image"
	ScreenshotModeFinviz      ScreenshotMode = "Finviz"
	ScreenshotModeMarketWatch ScreenshotMode = "MarketWatch"
	ScreenshotModeMarketBeat  ScreenshotMode = "MarketBeat"
	ScreenshotModeCathiesArk  ScreenshotMode = "CathiesArk"
	ScreenshotModeGuruFocus   ScreenshotMode = "GuruFocus"
	ScreenshotModeTipRanks    ScreenshotMode = "TipRanks"
	ScreenshotModeBarChart    ScreenshotMode = "BarChart"
)

// ArticleCase struct
type ArticleCase struct {
	name           string
	description    string
	linkURL        string
	imageURL       string
	screenshotMode ScreenshotMode
	// x              float64
	// y              float64
	// width          float64
	// height         float64
}

// TODO: чаще используемые перемещать наверх
// TODO: конфигурировать через интерфейс

// ArticleCases slice
var ArticleCases = []ArticleCase{
	{name: "tradingview", linkURL: "https://tradingview.com/symbols/%s"},
	{name: "chart", linkURL: "https://finviz.com/quote.ashx?t=%s",
		description:    "Chart",
		screenshotMode: ScreenshotModeImage,
		imageURL:       "https://charts2.finviz.com/chart.ashx?t=%s&ta=1&p=d&r=%d",
	},
	{name: "finviz", linkURL: "https://finviz.com/quote.ashx?t=%s",
		description:    "Financial Visualizations",
		screenshotMode: ScreenshotModeFinviz,
		// y:              215,
		// height:         845 - 91, /* (banner) */
	},
	{name: "gurufocus", linkURL: "https://gurufocus.com/stock/%s/summary#",
		description:    "Summary",
		screenshotMode: ScreenshotModeGuruFocus,
	},
	{name: "cathiesark", linkURL: "https://cathiesark.com/ark-combined-holdings-of-%s",
		description:    "ARK Invest Fund Holdings",
		screenshotMode: ScreenshotModeCathiesArk,
	},
	{name: "marketwatch", linkURL: "https://marketwatch.com/investing/stock/%s",
		description:    "Daily Price",
		screenshotMode: ScreenshotModeMarketWatch,
	}, // y: 345,
	// height: 565,
	{name: "marketbeat", linkURL: "https://marketbeat.com/stocks/%s",
		description:    "Insider Trades & Institutional Ownership",
		screenshotMode: ScreenshotModeMarketBeat,
	},
	{name: "stockcharts", linkURL: "https://stockcharts.com/h-sc/ui?s=%s",
		description:    "Technical Analysis",
		screenshotMode: ScreenshotModeImage,
		// imageURL:       "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=3&g=1&i=t8007589923c&r=%d", // portrain & 3 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=3&g=1&i=t0143501833c&r=%d", // portrain w/o zoom & 3 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=3&g=1&i=t8844882045c&r=%d", // landscape & 3 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=3&g=1&i=t2411346931c&r=%d", // landscape w/o zoom & 3 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=5&g=1&i=t8110014273c&r=%d", // portrain w/o zoom & 5 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=5&g=1&i=t7762146583c&r=%d", // landscape w/o zoom & 5 bar
		// imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=5&g=1&i=t6193398740c&r=%d", // portrain w/o zoom & 5 bar & default theme
		imageURL: "https://stockcharts.com/c-sc/sc?s=%s&p=D&b=5&g=1&i=t8674635596c&r=%d", // landscape w/o zoom & 5 bar & default theme
	},
	{name: "stockscores", linkURL: "https://stockscores.com/charts/charts/?ticker=%s",
		description:    "Technical Analysis",
		screenshotMode: ScreenshotModeImage,
		imageURL:       "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=180&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1180&ChartHeight=590&LogScale=None&Band=BB&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=None&Indicator3=MACD&Indicator4=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen&noCache=%d",
		// width:          1200,
		// height:         1000,
		// imageURL:       "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=120&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1200&ChartHeight=525&LogScale=None&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=100&Indicator1=RSI&Indicator2=None&Indicator3=MACD&Indicator4=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen&noCache=%d",
	},
	{name: "shortvolume", linkURL: "https://shortvolume.com/?t=%s",
		description:    "Daily Short Sale Volume",
		screenshotMode: ScreenshotModeImage,
		// width:          800,
		// height:         600,
		imageURL: "https://shortvolume.com/chart_engine/draw_chart.php?Symbol=%s&TimeRange=100&noCache=%d",
	},
	{name: "tipranks", linkURL: "https://tipranks.com/stocks/%s/stock-analysis",
		description: "Stock Analysis & Ratings",
		// screenshotMode: ScreenshotModePage,
		// x:              64,
		// y:              170,
		// width:          800 - 64,
		// height:         913,
		screenshotMode: ScreenshotModeTipRanks,
	},
	{name: "barchart", linkURL: "https://barchart.com/stocks/quotes/%s/overview",
		description:    "Overview",
		screenshotMode: ScreenshotModeBarChart,
	},
	{name: "finasquare", linkURL: "https://www.finasquare.com/stocks/%s/company-info/overview", description: "Overview"},
	{name: "stockrow", linkURL: "https://stockrow.com/%s", description: "Overview"},
	{name: "stockanalysis", linkURL: "https://stockanalysis.com/stocks/%s/", description: "Overview"},
	{name: "earningswhispers", linkURL: "https://earningswhispers.com/stocks/%s", description: "Overview"},
}

// TODO: articleCase by ScreenshotMode
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
