package chat

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var minWait = 4

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

func addPic(ctx context.Context, listOfPics []string, numberFrom string, numberTo string) {
	pathToPic, err := os.Getwd()
	if err != nil {
		log.Println("Error attach pic:", err)
		return
	}
	pathToPic = pathToPic + "\\pics\\" + randPic(listOfPics)

	err = chromedp.Run(ctx,
		chromedp.Click("span[data-icon='plus']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Println("Ошибка нажатия на кнопку прикрепления:", err)
		return
	}

	err = chromedp.Run(ctx,
		chromedp.SetUploadFiles("input[accept='image/*,video/mp4,video/3gpp,video/quicktime']", []string{pathToPic}, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		log.Println("Ошибка загрузки фото:", err)
		return
	}
}

func sendVoice(ctx context.Context) {
	wait := rand.Intn(25) + 1
	//log.Printf("Начало записи гс")

	err := chromedp.Run(ctx,
		chromedp.Click("button[aria-label='Голосовое сообщение']", chromedp.BySearch),
		chromedp.Sleep(time.Duration(wait)*time.Second),
	)
	if err != nil {
		log.Println("Error send voice:", err)
		return
	}

	//log.Printf("Тык на отправку гс")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("button[data-tab='11'][aria-label='Отправить']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),

		chromedp.Click("button[data-tab='11'][aria-label='Отправить']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Println("Error send voice:", err)
		return
	}

	//log.Printf("Жду изменения иконки гс")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("span[data-icon='ptt']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error send voice:", err)
		return
	}
	//log.Printf("Конец отправки гс")
}

func sendMessage(ctx context.Context, numberTo string, message string, listOfPics []string, needToWait bool, numberFrom string) error {

	log.Printf("отправляю %s > %s 2", numberFrom, numberTo)

	err := chromedp.Run(ctx,
		chromedp.Evaluate(`window.focus();`, nil),
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible("button[title='Новый чат']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
		chromedp.Click("button[title='Новый чат']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Println("Error create new chat:", err)
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[title='Новый чат']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys("div[class='x1hx0egp x6ikm8r x1odjw0f x6prxxf x1k6rcq7 x1whj5v'][tabindex='3'][data-tab='3']", numberTo, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error finding contact:", err)
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[class='_ak72 false _ak73 _ak7n'], div[class='_ak72 false _ak73']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.Click("div[class='_ak72 false _ak73 _ak7n'], div[class='_ak72 false _ak73']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error clicking contact:", err)
		return err
	}

	if needToWait {
		//log.Printf("Жду сообщения от %s чтобы отправить встречное", numberTo)
		err = waitForLastMessageIncoming(ctx)
		if err != nil {
			log.Printf("Ошибка ожидания сообщения: %v", err)
			return err
		}
		//log.Printf("Отправляю сообщение %s встречное получено", numberTo)
	}

	if rand.Float64() < 0.3 {
		sendVoice(ctx)
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible("div[aria-placeholder='Введите сообщение']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys("div[aria-placeholder='Введите сообщение']", message, chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}

	if rand.Float64() < 0.1 && len(listOfPics) != 0 {
		addPic(ctx, listOfPics, numberFrom, numberTo)
	}

	err = chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.WaitVisible("span[data-icon='send']", chromedp.BySearch),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Println("Error waiting for message input field:", err)
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.Click("span[data-icon='send']", chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible("span[data-icon='ptt']", chromedp.BySearch),
	)
	if err != nil {
		log.Println("Error sending message:", err)
		return err
	}

	log.Printf("Конец отправки сообщения %s > %s 2", numberFrom, numberTo)

	return nil
	//log.Printf("Сообщение отправлено номеру %s\n", numberTo, time.Now().Format("15:04:05"))
}

func waitForLastMessageIncoming(ctx context.Context) error {
	script := `
        (function() {
            const allMessages = document.querySelectorAll("div.message-in, div.message-out");
            if (allMessages.length === 0) {
                return false; // Нет сообщений
            }
            const lastMessage = allMessages[allMessages.length - 1];
            return lastMessage.classList.contains("message-in");
        })();
    `
	counter := 1

	for {
		var isIncoming bool
		err := chromedp.Run(ctx, chromedp.Evaluate(script, &isIncoming))
		if err != nil {
			return fmt.Errorf("ошибка выполнения JavaScript: %w", err)
		}

		if isIncoming {
			log.Println("Последнее сообщение входящее.")
			return nil
		}

		if counter == 2 {
			return fmt.Errorf("ошибка не дождался входящего")
		}
		time.Sleep(5 * time.Second)
		counter++
	}
}

func Chatting(parentCtx context.Context, syncch chan int, numberFrom string, numbersTo []string, duration time.Duration, n int, listOfPics []string) {

	val := <-syncch

	folderPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error get current dir: ", err, " for num: ", numberFrom)
	}
	folderPath = folderPath + "\\chromedp-profile\\" + numberFrom

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("user-data-dir", folderPath),
		chromedp.Flag("use-fake-ui-for-media-stream", true),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(parentCtx, opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.Navigate("https://web.whatsapp.com/"),
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal("Error opening WhatsApp Web:", err)
		return
	}

	log.Printf("Отсканируйте при надобности qr-код для номера %s\n", numberFrom)

	err = chromedp.Run(ctx,
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitVisible("span[data-icon='new-chat-outline']", chromedp.BySearch),
	)
	if err != nil {
		log.Println("Error create new chat:", err)
		return
	}

	for _, number := range numbersTo {
		if strings.Compare(number, numberFrom) == 0 {
			continue
		}
		err = sendMessage(ctx, number, randHello(), listOfPics, false, numberFrom)
		if err != nil {
			log.Println("Error:", err)
			return
		}
	}

	syncch <- (val + 1)

	maxSec := int(duration.Seconds()) / n
	maxSec = 5
	var wg sync.WaitGroup
	wg.Add(len(numbersTo))
	var mutex sync.Mutex

	log.Printf("Отправка сообщений с номера %s", numberFrom)
	log.Println(numbersTo)

	for _, number := range numbersTo {
		if strings.Compare(number, numberFrom) == 0 {
			continue
		}
		go func(to string) {
			defer wg.Done()
			i := 0
			for i <= n-2 {
				select {
				case <-ctx.Done():
					log.Printf("Остановка горутины %s > %s", numberFrom, to)
					return
				default:
					wait := rand.Intn(maxSec)
					time.Sleep(time.Duration(minWait+wait) * time.Second)

					var message string
					if i == n-2 {
						message = randBye()
					} else {
						message = randMessage()
					}

					log.Printf("отправляю %s - %s 1", numberFrom, number)
					mutex.Lock()
					err2 := sendMessage(ctx, to, message, listOfPics, true, numberFrom)
					mutex.Unlock()

					if err2 != nil {
						log.Println("Error:", err2)
						continue
					}

					i++
				}
			}
		}(number)
	}
	wg.Wait()

	log.Printf("Все сообщения от %s отправлены\n", numberFrom)
}
