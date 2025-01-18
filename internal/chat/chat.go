package chat

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func sendMessage(ctx context.Context, numberTo string, message string) {

	err := chromedp.Run(ctx,
		chromedp.Click("span[data-icon='new-chat-outline']", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error click new chat:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[contenteditable='true'][data-tab='3']", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error finding contact:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.SendKeys("div[contenteditable='true'][data-tab='3']", numberTo, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println("Error searching contact:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible(fmt.Sprintf("span[title='+7 %s'", numberTo), chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error finding contact:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.Click("div[tabindex='0'][role='button']", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error clicking contact:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("//p[contains(@class, 'selectable-text') and contains(@class, 'copyable-text') and contains(@class, 'x15bjb6t') and contains(@class, 'x1n2onr6')]",
			chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error waiting for message input field:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.SendKeys("//*[@id='main']/footer/div[1]/div/span/div/div[2]/div[1]/div[2]/div[1]", message, chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.Click("//*[@id='main']/footer/div[1]/div/span/div/div[2]/div[2]/button/span", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	log.Println("Сообщение отправлено номеру %s в %s", numberTo, time.Now().Format("15:04:05"))
}

func Chatting(syncch chan interface{}, numberFrom string, numbers map[string]int, n int, isWarm int) {

	<-syncch

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com/"),
	)
	if err != nil {
		log.Fatal("Error opening WhatsApp Web:", err)
	}

	log.Println("Отсканируйте qr-код для номера %s", numberFrom)

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("span[data-icon='new-chat-outline']", chromedp.BySearch),
	)
	if err != nil {
		fmt.Println("Error create new chat:", err)
		return
	}

	syncch <- struct{}{}

	numbersTo := make([]string, 0, len(numbers))
	if isWarm == 0 {
		for k := range numbers {
			if strings.Compare(k, numberFrom) != 0 {
				numbersTo = append(numbersTo, k)
			}
		}

	} else {
		for k, v := range numbers {
			if strings.Compare(k, numberFrom) != 0 && v != 1 {
				numbersTo = append(numbersTo, k)
			}
		}
	}

	maxSec := 1 * 5 * 60 / n
	var wg sync.WaitGroup
	wg.Add(len(numbersTo))
	var mutex sync.Mutex
	log.Println(numbersTo)

	for _, number := range numbersTo {
		go func() {
			for _ = range n {
				wait := rand.Intn(maxSec)
				mutex.Lock()
				sendMessage(ctx, number, "bbbbb")
				mutex.Unlock()
				time.Sleep(time.Duration(wait) * time.Second)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("Все сообщения от %s отправлены", numberFrom)
}
