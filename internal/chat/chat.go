package chat

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func randHello() string {
	return randomChoice(hi)
}

func randBye() string {
	return randomChoice(bye)
}

func randMessage() (message string) {

	numPhrases := rand.Intn(4) + 1

	messageParts := randomChoices(phrases, numPhrases)

	// Вставляем случайное слово в начало, середину или конец
	if rand.Float64() > 0.5 {
		insertIndex := rand.Intn(len(messageParts) + 1)
		messageParts = append(messageParts[:insertIndex], append([]string{randomChoice(inserts)}, messageParts[insertIndex:]...)...)
	}

	// Собираем сообщение
	message = strings.Join(messageParts, " ")

	// Добавляем смайлики с вероятностью
	if rand.Float64() > 0.5 {
		message += " " + randomChoice(emojis)
	}
	if rand.Float64() > 0.95 {
		message += " " + randomChoice(emojis)
	}

	return message
}

func randomChoice(options []string) string {
	return options[rand.Intn(len(options))]
}

func randomChoices(options []string, n int) []string {
	var results []string
	for i := 0; i < n; i++ {
		results = append(results, randomChoice(options))
	}
	return results
}

func randPic(listOfPics []string) string {
	return randomChoice(listOfPics)
}

func addPic(ctx context.Context, listOfPics []string) {
	pathToPic := "..\\pics\\" + randPic(listOfPics)

	err := chromedp.Run(ctx,
		chromedp.Click("span[data-icon='plus']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),

		chromedp.SetUploadFiles("input[accept='image/*,video/mp4,video/3gpp,video/quicktime'][type='file']", []string{pathToPic}, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Println("Error attach pic:", err)
		return
	}
}

func sendVoice(ctx context.Context) {
	wait := rand.Intn(25) + 1

	err := chromedp.Run(ctx,
		chromedp.Click("button[aria-label='Голосовое сообщение']", chromedp.BySearch),
		chromedp.Sleep(time.Duration(wait)*time.Second),

		chromedp.Click("button[data-tab='11'][aria-label='Отправить']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),

		chromedp.WaitVisible("span[data-icon='ptt']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error send voice:", err)
		return
	}
}

func sendMessage(ctx context.Context, numberTo string, message string, listOfPics []string) {

	err := chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible("button[title='Новый чат']", chromedp.BySearch),
		chromedp.Click("button[title='Новый чат']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Println("Error create new chat:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[title='Новый чат']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys("div[class='x1hx0egp x6ikm8r x1odjw0f x6prxxf x1k6rcq7 x1whj5v'][tabindex='3'][data-tab='3']", numberTo, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error finding contact:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[class='_ak72 false _ak73 _ak7n'], div[class='_ak72 false _ak73']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.Click("div[class='_ak72 false _ak73 _ak7n'], div[class='_ak72 false _ak73']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error clicking contact:", err)
		return
	}

	if rand.Float64() < 0.3 {
		sendVoice(ctx)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[data-tab='10']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys("div[data-tab='10']", message, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error writing message:", err)
		return
	}

	if rand.Float64() < 0.3 && len(listOfPics) != 0 {
		addPic(ctx, listOfPics)
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("span[data-icon='send']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error waiting for message input field:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.Click("button[data-tab='11']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	//log.Printf("Сообщение отправлено номеру %s\n", numberTo, time.Now().Format("15:04:05"))
}

func Chatting(syncch chan interface{}, numberFrom string, numbers map[string]int, duration time.Duration, n int, listOfPics []string, isWarm int) {

	<-syncch

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("use-fake-ui-for-media-stream", true),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com/"),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal("Error opening WhatsApp Web:", err)
	}

	log.Printf("Отсканируйте qr-код для номера %s\n", numberFrom)

	err = chromedp.Run(ctx,
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitVisible("span[data-icon='new-chat-outline']", chromedp.BySearch),
	)
	if err != nil {
		log.Println("Error create new chat:", err)
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

	maxSec := int(duration.Seconds()) / n
	var wg sync.WaitGroup
	wg.Add(len(numbersTo))
	var mutex sync.Mutex
	var message string
	//log.Println(numbersTo)

	for _, number := range numbersTo {
		go func() {
			for i := range n {
				wait := rand.Intn(maxSec)
				if i == 0 {
					message = randHello()
				} else if i == n-1 {
					message = randBye()
				} else {
					message = randMessage()
				}
				mutex.Lock()
				sendMessage(ctx, number, message, listOfPics)
				mutex.Unlock()

				time.Sleep(time.Duration(wait+1) * time.Second)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	log.Printf("Все сообщения от %s отправлены\n", numberFrom)
}
