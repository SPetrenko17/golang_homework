package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func readStrings(file *os.File) ([]string, int, error) {
	regSpace := regexp.MustCompile(" ")
	minColumns := math.MaxInt64
	i := 0
	var res = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		res = append(res, scanner.Text())
		if len(regSpace.Split(res[len(res)-1], -1)) < minColumns {
			minColumns = len(strings.Fields(res[len(res)-1]))
		}
	}
	if err := scanner.Err(); err != nil {
		return res, minColumns, err
	}
	return res, minColumns, nil
}

func isNumbersSlice(strSlice []string) bool {
	regexWrongSymbols := regexp.MustCompile("[^0-9 ]")
	for i := 0; i < len(strSlice); i++ {
		wrongSymbols := regexWrongSymbols.FindAllString(strSlice[i], -1)
		if len(wrongSymbols) > 0 {
			return false
		}
	}
	return true

}

func validate(sortByColumn int, minColumns int, sortByNumbers bool, stringsFromFile []string) error {
	if sortByColumn != 0 && sortByColumn > minColumns {
		return errors.New("invaild column count ")
	}

	if sortByNumbers && !isNumbersSlice(stringsFromFile) {
		return errors.New("invaild number file")
	}

	return nil
}

func getFile() {

}
func main() {
	ignoreUppercase := flag.Bool("f", false, "Ignore uppercase")
	uniqueValues := flag.Bool("u", false, "Unique values")
	sortDescending := flag.Bool("r", false, "Sort Descending")
	sortByNumbers := flag.Bool("n", false, "Sort by numbers")
	fileOutput := flag.String("o", "", "Output in file")
	sortByColumn := flag.Int("k", 0, "Sort by column (k-word in string)")
	flag.Parse()

	args := os.Args
	file, err := os.Open(args[len(args)-1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stringsFromFile, minColumns, err := readStrings(file)
	if err != nil {
		fmt.Println("error ReadStrings")
		os.Exit(1)
	}

	validationError := validate(*sortByColumn, minColumns, *sortByNumbers, stringsFromFile)
	if validationError != nil {
		os.Exit(1)
	}
	err = mySort(&stringsFromFile, *ignoreUppercase, *sortDescending, *sortByNumbers, *sortByColumn)
	if err != nil {
		os.Exit(1)
	}
	err = output(stringsFromFile, *fileOutput, *ignoreUppercase, *uniqueValues)
	if err != nil {
		os.Exit(1)
	}

}

func applyColumnSorting(left string, right string, columnsCount int) (string, string) {
	left = strings.Fields(left)[columnsCount-1]
	right = strings.Fields(right)[columnsCount-1]
	return left, right
}

func applyIgnoreUppercase(left string, right string) (string, string) {
	left = strings.ToLower(left)
	right = strings.ToLower(right)
	return left, right
}

func applySortByNumbers(left string, right string) (float64, float64, error) {
	leftInt, err := strconv.Atoi(left)
	if err != nil {
		return 0, 0, err
	}
	rightInt, err := strconv.Atoi(right)
	if err != nil {
		return 0, 0, err
	}
	return float64(leftInt), float64(rightInt), nil
}

func compare(left, right interface{}, isNumeric, IsDescending bool) bool {
	if isNumeric {
		if IsDescending {
			return left.(float64) > right.(float64)
		} else {
			return left.(float64) < right.(float64)
		}
	} else {
		if IsDescending {
			return left.(string) > right.(string)
		} else {
			return left.(string) < right.(string)
		}
	}
}

func mySort(strSlice *[]string, ignoreUppercase, sortDescending, sortByNumbers bool, sortByColumn int) error {
	sort.Slice(*strSlice, func(i, j int) bool {
		var left = (*strSlice)[i]
		var right = (*strSlice)[j]

		if sortByColumn > 0 {
			left, right = applyColumnSorting(left, right, sortByColumn)
		}
		if ignoreUppercase {
			left, right = applyIgnoreUppercase(left, right)
		}

		if sortByNumbers {
			leftInt, rightInt, err := applySortByNumbers(left, right)
			if err != nil {
				os.Exit(1)
			}
			return compare(leftInt, rightInt, true, sortDescending)
		} else {
			return compare(left, right, false, sortDescending)
		}
	})
	return nil
}

func output(strings []string, filename string, uniqueValues, ignoreUppercase bool) error {
	var f *os.File
	var err error
	var isFileOutput = true

	if filename != "" {
		f, err = os.Create(filename)
		if err != nil {
			isFileOutput = false
		}
	} else {
		isFileOutput = false
	}

	for i := 0; i < len(strings); i++ {
		left := strings[i]
		right := ""
		if i != len(strings)-1 {
			right = strings[i+1]
		}
		if ignoreUppercase {
			left, right = applyIgnoreUppercase(left, right)
		}
		if uniqueValues && left == right {
			i++
			if isFileOutput {
				_, err := fmt.Fprintln(f, strings[i])
				if err != nil {
					return err
				}
			} else {
				fmt.Println(strings[i])
			}

		} else {
			if isFileOutput {
				_, err := fmt.Fprintln(f, strings[i])
				if err != nil {
					return err
				}
			} else {
				fmt.Println(strings[i])
			}
		}
	}
	return nil
}
