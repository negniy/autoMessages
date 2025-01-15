package main

import (
	"autoMessages/internal/config"
	"context"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func run(numbers map[string]int) {

	var wg sync.WaitGroup
	wg.Add(len(numbers))

	for number, isWarm := range len(numbers) {
		go func() {
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
			time.Sleep(30 * time.Second)
			wg.Done()
		}()
	}

	wg.Wait()
	log.Println("Номера прогреты")

}

func initNumbers() map[string]int {
	numbers := config.LoadNumbers()
	return numbers
}

func main() {

}
