package main

import (
	"log"
	"time"
)

// TODO: подключить ETF-ки https://etfdb.com/screener/
// TODO: выдавать сообщение sendLink, а по готовности основного ответа - его удалять
// TODO: кнопки под полем ввода в приватном чате для: inline mode, help, search & all,
// TODO: реализовать румтур
// TODO: поиск по ticker.title
// TODO: README
// TODO: svg to png
// TODO: добавить тайм-фрейм #BABA?15M
// TODO: добавить медленную скользящую #BABA?50EMA / 100EMA / 200EMA
// TODO: параллельная обрарботка https://gobyexample.ru/worker-pools.html
// TODO: добавить ETF, например ARKK
// TODO: добавить биток GBTC
// TODO: добавить опционы с investing.com
// TODO: не успевает загрузить картинку tipranks.com (показывает колёсики)

func main() {
	go backgroundTask()
	// // This print statement will be executed before
	// // the first `tock` prints in the console
	log.Println("The rest of my application can continue")
	// // here we use an empty select{} in order to keep
	// // our main function alive indefinitely as it would
	// // complete before our backgroundTask has a chance
	// // to execute if we didn't.
	select {}
}

func backgroundTask() {
	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		log.Println("Tick at", t.Minute(), t.Minute()%10, t.Second())
		// t.Minute()%10 == 0 &&
		if t.Second() == 4 {
			// if sendFinvizMap(b, chatID) {
			log.Println("Send map")
			// }
		}
	}
}
