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
		chromedp.WaitVisible(fmt.Sprintf("//span[@title='%s']", numberTo), chromedp.BySearch),
	)
	if err != nil {
		log.Fatal("Error finding contact:", err)
	}
}

func Chatting(numberFrom string, numbers map[string]int, n int, isWarm int) {

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

	maxSec := 8 * 60 * 60 / n
	var wg sync.WaitGroup
	wg.Add(len(numbersTo))
	var mutex sync.Mutex

	for _, number := range numbersTo {
		go func() {
			for _ = range n {
				wait := rand.Intn(maxSec)
				mutex.Lock()
				sendMessage(ctx, number, "aboba")
				mutex.Unlock()
				time.Sleep(time.Duration(wait) * time.Second)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("Все сообщения от %s отправлены", numberFrom)
}
