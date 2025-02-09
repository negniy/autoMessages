package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CreateFolder() {}

func LoadPics(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	reg, err := regexp.Compile(`.*\.jpg$`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pics := make([]string, 0, 10)
	for _, file := range files {
		if !file.IsDir() {
			if reg.MatchString(file.Name()) {
				pics = append(pics, file.Name())
			}
		}
	}
	return pics, nil
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

func LoadNumbers(filename string) ([]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	numbers := make([]string, 0, 10)
	index := 1
	re, err := regexp.Compile(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`)
	if err != nil {
		log.Println(err)
		return nil, err
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
					return numbers, nil
				}
			}
			continue
		}

		number = normalizeNumber(number)
		number = formatPhoneNumber(number)
		numbers = append(numbers, number)
		index++
	}
}
