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

func run(warm []string, cold []string, duration time.Duration, n int, listOfPics []string) {

	var wg sync.WaitGroup
	wg.Add(len(warm) + len(cold))
	ch := make(chan interface{}, 1)
	ch <- struct{}{}

	for _, number := range warm {
		go func() {
			chat.Chatting(ch, number, cold, duration, n, listOfPics)
			wg.Done()
		}()
	}

	numbers := append(cold, warm...)

	for _, number := range cold {
		go func() {
			chat.Chatting(ch, number, numbers, duration, n, listOfPics)
			wg.Done()
		}()
	}

	wg.Wait()
	close(ch)
	log.Println("Номера прогреты")

}

func initWarmNumbers(filename string) ([]string, error) {
	numbers, err := config.LoadNumbers(filename)
	return numbers, err
}

func initColdNumbers(filename string) ([]string, error) {
	numbers, err := config.LoadNumbers(filename)
	return numbers, err
}

func parseDuration(input string) (time.Duration, error) {
	switch input {
	case "1.5 часа":
		return 1*time.Hour + 30*time.Minute, nil
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

	var warm []string
	var cold []string
	numberOfMessages := 0
	var err error
	var time time.Duration = 0
	var pics []string

	logWriter := &logger.CustomWriter{}
	log.SetOutput(logWriter)

	myApp := app.New()
	myWindow := myApp.NewWindow("Прогрев")

	logOutput := widget.NewMultiLineEntry()
	logOutput.SetPlaceHolder("Здесь будут отображаться логи...")
	logOutput.Disable()
	logWriter.SetLogWidget(logOutput)

	messageCountButton := widget.NewButton("Выбор количества сообщений", func() {
		options := []string{"10", "20", "30", "40", "50", "60", "70", "80", "90", "100"}
		dialog.ShowCustom("Выберите количество сообщений", "OK", container.NewVBox(
			widget.NewLabel("Количество сообщений:"),
			widget.NewSelect(options, func(selected string) {
				log.Printf("Выбрано количество сообщений: %s\n", selected)
				numberOfMessages, err = strconv.Atoi(selected)
				if err != nil {
					log.Printf("Ошибка: %v\n", err)
					return
				}
			}),
		), myWindow)
	})

	timeSelectionButton := widget.NewButton("Выбор времени работы", func() {
		options := []string{"1.5 часа", "2 часа", "3 часа", "5 часов", "8 часов", "12 часов"}
		dialog.ShowCustom("Выберите время работы", "OK", container.NewVBox(
			widget.NewLabel("Время работы:"),
			widget.NewSelect(options, func(selected string) {
				log.Printf("Выбрано время работы: %s\n", selected)
				time, err = parseDuration(selected)
				if err != nil {
					log.Printf("Ошибка: %v\n", err)
					return
				}
			}),
		), myWindow)
	})

	pics, err = config.LoadPics("./pics")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		return
	}
	log.Println(pics)

	warm, err = config.LoadNumbers(".//dontTouch.xlsx")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		return
	}
	log.Println("Warm: ", warm)

	cold, err = config.LoadNumbers(".//numbers.xlsx")
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		return
	}
	log.Println("Cold: ", cold)

	startButton := widget.NewButton("Запустить прогрев", func() {
		if len(warm) == 0 || len(cold) == 0 {
			log.Println("При загрузке номеров возникла ошибка")
			return
		} else if time == 0 {
			log.Println("При выборе времени прогрева возникла ошибка")
			return
		} else if numberOfMessages == 0 {
			log.Println("При выборе кол-ва сообщений возникла ошибка")
			return
		} else if len(pics) == 0 {
			log.Println("При загрузке картинок возникла ошибка")
			return
		} else {
			run(warm, cold, time, numberOfMessages, pics)
		}
	})

	content := container.NewVBox(
		messageCountButton,
		timeSelectionButton,
		startButton,
		logOutput,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.ShowAndRun()

}
