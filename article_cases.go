package main

// ArticleCase struct
type ArticleCase struct {
	name    string
	url     string
	preview bool
}

// ArticleCases slice
var ArticleCases = []ArticleCase{
	{name: "finviz.com", url: "https://finviz.com/quote.ashx?t=%s"},
	{name: "cathiesark.com", url: "https://cathiesark.com/ark-combined-holdings-of-%s", preview: true},
	{name: "marketwatch.com", url: "https://marketwatch.com/investing/stock/%s"},
	{name: "nakedshortreport.com", url: "https://nakedshortreport.com/company/%s"},
	{name: "gurufocus.com", url: "https://gurufocus.com/stock/%s/summary"},
	// {name: "stockanalysis.com", url: "https://stockanalysis.com/stocks/%s/financials/"},
	// {name: "stockrow.com", url: "https://stockrow.com/%s"},
}
