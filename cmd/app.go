package main

import (
	"autoMessages/internal/chat"
	"autoMessages/internal/config"
	"autoMessages/internal/logger"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func run(numbers map[string]int, duration time.Duration, n int) {

	var wg sync.WaitGroup
	wg.Add(len(numbers))
	ch := make(chan interface{}, 1)
	ch <- struct{}{}

	warm := make([]string, 0, len(numbers))
	cold := make([]string, 0, len(numbers))
	for k, v := range numbers {
		if v == 1 {
			warm = append(warm, k)
		} else {
			cold = append(cold, k)
		}
	}

	for _, number := range warm {
		go func() {
			chat.Chatting(ch, number, numbers, duration, n, 1)
			wg.Done()
		}()
	}

	for _, number := range cold {
		go func() {
			chat.Chatting(ch, number, numbers, duration, n, 0)
			wg.Done()
		}()
	}

	wg.Wait()
	close(ch)
	log.Println("Номера прогреты")

}

func initNumbers(filename string) (map[string]int, error) {
	numbers, err := config.LoadNumbers(filename)
	return numbers, err
}

func parseDuration(input string) (time.Duration, error) {
	switch input {
	case "5 мин":
		return 5 * time.Minute, nil
	case "10 мин":
		return 10 * time.Minute, nil
	case "30 мин":
		return 30 * time.Minute, nil
	case "1 час":
		return 1 * time.Hour, nil
	case "2 часа":
		return 2 * time.Hour, nil
	case "3 часа":
		return 3 * time.Hour, nil
	case "5 часов":
		return 5 * time.Hour, nil
	case "8 часов":
		return 8 * time.Hour, nil
	case "12 часов":
		return 12 * time.Hour, nil
	default:
		return 0, fmt.Errorf("неизвестное время: %s", input)
	}
}

func main() {

	var numbers map[string]int
	numberOfMessages := 0
	var err error
	var time time.Duration = 0

	logWriter := &logger.CustomWriter{}
	logWriter.AddOutput(log.Writer())
	log.SetOutput(logWriter)

	myApp := app.New()
	myWindow := myApp.NewWindow("Прогрев")

	logOutput := widget.NewMultiLineEntry()
	logOutput.SetPlaceHolder("Здесь будут отображаться логи...")
	logOutput.Disable()
	logWriter.SetLogWidget(logOutput)

	config.LoadPics()
	log.Println("List of ", chat.Pics)

	messageCountButton := widget.NewButton("Выбор количества сообщений", func() {
		options := []string{"10", "20", "30", "40", "50", "60", "70", "80", "90", "100"}
		dialog.ShowCustom("Выберите количество сообщений", "OK", container.NewVBox(
			widget.NewLabel("Количество сообщений:"),
			widget.NewSelect(options, func(selected string) {
				log.Println(fmt.Sprintf("Выбрано количество сообщений: %s", selected))
				numberOfMessages, err = strconv.Atoi(selected)
				if err != nil {
					log.Println(fmt.Sprintf("Ошибка: %v", err))
					return
				}
			}),
		), myWindow)
	})

	timeSelectionButton := widget.NewButton("Выбор времени работы", func() {
		options := []string{"5 мин", "10 мин", "30 мин", "1 час", "2 часа", "3 часа", "5 часов", "8 часов", "12 часов"}
		dialog.ShowCustom("Выберите время работы", "OK", container.NewVBox(
			widget.NewLabel("Время работы:"),
			widget.NewSelect(options, func(selected string) {
				log.Println(fmt.Sprintf("Выбрано время работы: %s", selected))
				time, err = parseDuration(selected)
				if err != nil {
					log.Println(fmt.Sprintf("Ошибка: %v", err))
					return
				}
			}),
		), myWindow)
	})

	excelButton := widget.NewButton("Открыть Excel с номерами", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println(fmt.Sprintf("Ошибка: %v", err))
				return
			}
			if reader == nil {
				log.Println("Файл не выбран")
				return
			}
			fileName := reader.URI().String()[7:]
			log.Println(fileName)
			numbers, err = initNumbers(fileName)
			if err != nil {
				log.Println(fmt.Sprintf("Ошибка: %v", err))
				return
			}
			log.Println(fmt.Sprintf("Выбран файл: %s", fileName))
		}, myWindow).Show()
	})

	startButton := widget.NewButton("Запустить прогрев", func() {
		log.Println("Запуск прогрева...")
		if len(numbers) == 0 || time == 0 || numberOfMessages == 0 {
			log.Println("При подготовке возникла ошибка")
			return
		}
		run(numbers, time, numberOfMessages)
	})

	content := container.NewVBox(
		messageCountButton,
		timeSelectionButton,
		excelButton,
		startButton,
		logOutput,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.ShowAndRun()

}
