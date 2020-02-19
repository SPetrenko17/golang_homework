package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
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
		return
	}

	validationError := validate(*sortByColumn, minColumns, *sortByNumbers, stringsFromFile)
	if validationError != nil {
		fmt.Println("error validation")
		return
	}
	err = mySort(&stringsFromFile, *ignoreUppercase, *sortDescending, *sortByNumbers, *sortByColumn)
	if err != nil {
		fmt.Println("error mySort")
		return
	}
	err = output(stringsFromFile, *fileOutput, *uniqueValues, *ignoreUppercase)
	if err != nil {
		fmt.Println("error output")
		return
	}

}

func applyColumnSorting(left string, right string, columnsCount int) (leftColumn, rightColumn string) {
	leftColumn = strings.Fields(left)[columnsCount-1]
	rightColumn = strings.Fields(right)[columnsCount-1]
	return
}

func applyIgnoreUppercase(left string, right string) (leftLower, rightLower string) {
	leftLower = strings.ToLower(left)
	rightLower = strings.ToLower(right)
	return
}

func applySortByNumbers(left string, right string) (leftFloat, rightFloat float64) {
	leftInt, _ := strconv.Atoi(left)
	rightInt, _ := strconv.Atoi(right)
	leftFloat = float64(leftInt)
	rightFloat = float64(rightInt)
	return
}

func compare(left, right interface{}, IsDescending bool) bool {
	lt := reflect.TypeOf(left).Kind()
	rt := reflect.TypeOf(right).Kind()
	if lt == reflect.Float64 && rt == reflect.Float64 {
		if IsDescending {
			return left.(float64) > right.(float64)
		}
		return left.(float64) < right.(float64)
	}
	if IsDescending {
		return left.(string) > right.(string)
	}
	return left.(string) < right.(string)

}

func validateNumbersAndColumns(strSlice []string, sortByNumbers bool, sortByColumn int) error {
	if sortByNumbers && sortByColumn == 0 {
		isNumSl := isNumbersSlice(strSlice)
		if !isNumSl {
			return errors.New("not number slice")
		}
	} else if sortByNumbers && sortByColumn != 0 {
		column := make([]string, 0)
		for i := 0; i < len(strSlice); i++ {
			column = append(column)
		}
		isNumSl := isNumbersSlice(column)
		if !isNumSl {
			return errors.New("not number slice")
		}
	}
	return nil
}

func mySort(strSlice *[]string, ignoreUppercase, sortDescending, sortByNumbers bool, sortByColumn int) error {
	err := validateNumbersAndColumns(*strSlice, sortByNumbers, sortByColumn)
	if err != nil {
		return errors.New("validation error ")
	}
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
			leftInt, rightInt := applySortByNumbers(left, right)
			return compare(leftInt, rightInt, sortDescending)
		}
		return compare(left, right, sortDescending)
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
			return errors.New("error read file")
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
				continue
			}
			fmt.Println(strings[i])
			continue
		}
		if isFileOutput {
			_, err := fmt.Fprintln(f, strings[i])
			if err != nil {
				return err
			}
			continue
		}
		fmt.Println(strings[i])

	}
	return nil
}
