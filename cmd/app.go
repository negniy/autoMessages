package main

import (
	"autoMessages/internal/chat"
	"autoMessages/internal/config"
	"log"
	"sync"
)

func run(numbers map[string]int) {

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
			chat.Chatting(ch, number, numbers, 5, 1)
			wg.Done()
		}()
	}

	for _, number := range cold {
		go func() {
			chat.Chatting(ch, number, numbers, 5, 0)
			wg.Done()
		}()
	}

	wg.Wait()
	close(ch)
	log.Println("Номера прогреты")

}

func initNumbers() map[string]int {
	numbers := config.LoadNumbers()
	return numbers
}

func main() {

	numbers := initNumbers()
	log.Println("List of numbers", numbers)

	config.LoadPics()
	log.Println("List of ", chat.Pics)

}
