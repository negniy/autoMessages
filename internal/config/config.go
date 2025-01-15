package config

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func LoadNumbers() map[string]int {
	f, err := excelize.OpenFile("numbers.xlsx")
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
		numbers[number] = flag
		index++
	}
}
