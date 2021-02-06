package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gosimple/slug"
)

func main() {
	// Стартуем хром
	ctx, cancel := chromedp.NewContext(context.Background())
	// Не забываем, что его надо закрыть
	// при выходе из main вызовется cancel и он передаст
	// хрому о закрытии
	defer cancel()

	// первым аргументом скрипт принимает адрес страницы
	url := os.Args[1]

	// второй, опциональный аргумент, куда записывать скриншот
	// если не передан, то сохраним в png файл, соответствующий url
	// вырезав все лишнее из него библиотечкой slug
	var filename string
	if len(os.Args) == 3 {
		filename = os.Args[2]
	} else {
		filename = slug.Make(url) + ".png"
	}

	// инициализируем пустой массив, куда будет сохранен скриншот
	var imageBuf []byte

	// и отправляем хрому задачи, которые он должен выполнить
	// у нас только одна - ScreenshotTsks, но можно закинуть сколько угодно
	if err := chromedp.Run(
		ctx,
		ScreenshotTasks(url, &imageBuf),
	); err != nil {
		log.Fatal(err)
	}

	// Задача выполнена, можно сохранить полученное изображение в файл
	if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}
}

// ScreenshotTasks записывает в imageBuf скриншот страницы, расположенной на url
func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		// задача (таска) состоит из последовательности действий
		// сначала мы переходим по заданному url
		chromedp.Navigate(url),
		// а теперь делаем скриншот, записывая его в imageBuf
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().Do(ctx)
			return err
		}),
	}
}
