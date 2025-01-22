package config

import (
	"autoMessages/internal/chat"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func LoadPics() {
	files, err := os.ReadDir("..\\pics")
	if err != nil {
		log.Println(err)
	}
	reg, err := regexp.Compile(`.*\.jpg$`)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			if reg.MatchString(file.Name()) {
				chat.Pics = append(chat.Pics, file.Name())
			}
		}
	}
}

func normalizeNumber(number string) string {
	re := regexp.MustCompile(`^(?:\+7|7|8)`)
	return re.ReplaceAllString(number, "")
}

func formatPhoneNumber(input string) string {
	if len(input) != 10 {
		log.Printf("invalid phone number length")
		return ""
	}

	return fmt.Sprintf("%s %s-%s-%s",
		input[:3],  // Первая группа (999)
		input[3:6], // Вторая группа (123)
		input[6:8], // Третья группа (45)
		input[8:],  // Четвёртая группа (67)
	)
}

func LoadNumbers() map[string]int {
	f, err := excelize.OpenFile("../numbers.xlsx")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	numbers := make(map[string]int)
	index := 1
	re, err := regexp.Compile(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	for {
		number, err := f.GetCellValue("Лист1", "A"+strconv.Itoa(index))
		if err != nil {
			log.Println(err)
			continue
		}

		ok := re.MatchString(number)
		if !ok {
			log.Println("invalid number in " + "A" + strconv.Itoa(index))
			if strings.Compare(number, "") == 0 {
				if index != 1 {
					log.Println("reading is over at " + "A" + strconv.Itoa(index))
					return numbers
				}
			}
			continue
		}

		val, err := f.GetCellValue("Лист1", "B"+strconv.Itoa(index))
		if err != nil {
			log.Println(err)
			continue
		}
		flag, err := strconv.Atoi(val)
		if err != nil {
			log.Println(err)
			continue
		}

		if flag < 1 {
			flag = 0
		} else {
			flag = 1
		}

		number = normalizeNumber(number)
		number = formatPhoneNumber(number)
		numbers[number] = flag
		index++
	}
}
