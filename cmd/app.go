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

	for number, isWarm := range numbers {
		go func() {
			chat.Chatting(ch, number, numbers, 5, isWarm)
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
	//log.Println(numbers)
	run(numbers)

}
