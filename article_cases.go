package main

// ArticleCase struct
type ArticleCase struct {
	name       string
	url        string
	hasPreview bool
}

// TODO: чаще используемые перемещать наверх

// ArticleCases slice
var ArticleCases = []ArticleCase{
	{name: "finviz.com", url: "https://finviz.com/quote.ashx?t=%s"},
	{name: "cathiesark.com", url: "https://cathiesark.com/ark-combined-holdings-of-%s", hasPreview: true},
	{name: "stockscores.com", url: "https://www.stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=100&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=860&ChartHeight=480&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=50&Indicator1=Momentum&Indicator2=RSI&Indicator3=MACD&Indicator4=AccDist&endDate=2021-2-13&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		// "https://stockscores.com/chart.asp?TickerSymbol=%s&TimeRange=365&Interval=d&Volume=1&ChartType=CandleStick&Stockscores=None&ChartWidth=1920&ChartHeight=430&LogScale=1&Band=None&avgType1=EMA&movAvg1=20&avgType2=EMA&movAvg2=50&Indicator1=Momentum&Indicator2=RSI&Indicator3=MACD&Indicator4=AccDist&CompareWith=&entryPrice=&stopLossPrice=&candles=redgreen",
		hasPreview: true},
	{name: "shortvolume.com", url: "https://shortvolume.com/chart_engine/draw_chart.php?Symbol=%s&TimeRange=100", hasPreview: true},
	{name: "earningswhispers.com", url: "https://earningswhispers.com/stocks/%s"},
	{name: "marketwatch.com", url: "https://marketwatch.com/investing/stock/%s"},
	{name: "tipranks.com", url: "https://tipranks.com/stocks/%s/forecast"},
	{name: "barchart.com", url: "https://barchart.com/stocks/quotes/%s/overview"},
	{name: "gurufocus.com", url: "https://gurufocus.com/stock/%s/summary"},
	{name: "stockrow.com", url: "https://stockrow.com/%s"},
	{name: "finasquare.com", url: "https://www.finasquare.com/stocks/%s/company-info/overview"},
	{name: "marketbeat.com", url: "https://marketbeat.com/stocks/%s"},
	{name: "tradingview.com", url: "https://ru.tradingview.com/symbols/%s"},
}
