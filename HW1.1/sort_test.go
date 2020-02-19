package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestIsNumberSlice(t *testing.T) {
	var sl0 = make([]string, 0)
	sl0 = append(sl0, "1 2 3 4")
	sl0 = append(sl0, "5 6 7 8")
	sl0 = append(sl0, "9 10 11")
	assert.Equal(t, true, isNumbersSlice(sl0))

	var sl1 = make([]string, 0)
	sl1 = append(sl1, "1 2 3 4")
	sl1 = append(sl1, "5 6 7 8")
	sl1 = append(sl1, "hello moto")
	assert.Equal(t, false, isNumbersSlice(sl1))

	var sl2 = make([]string, 0)
	sl2 = append(sl2, "1 2 3 4")
	sl2 = append(sl2, "5 6 7 8")
	sl2 = append(sl2, "hello")
	assert.Equal(t, false, isNumbersSlice(sl1))

	var sl3 = make([]string, 0)
	sl3 = append(sl3, "1 2 3 4")
	sl3 = append(sl3, "5 6 7 8")
	sl3 = append(sl3, "9")
	assert.Equal(t, true, isNumbersSlice(sl0))
}

func TestCompare(t *testing.T) {
	assert.Equal(t, false, compare(1.0, 2.0, true))
	assert.Equal(t, true, compare(1.0, 2.0, false))
	assert.Equal(t, false, compare("abc", "def", true))
	assert.Equal(t, true, compare("abc", "def", false))
}

func TestApplyIgnoreUppercase(t *testing.T) {
	str1, str2 := applyIgnoreUppercase("HELLO", "MoTo")
	assert.Equal(t, str1, "hello")
	assert.Equal(t, str2, "moto")
	str1, str2 = applyIgnoreUppercase("hello", "moto")
	assert.Equal(t, str1, "hello")
	assert.Equal(t, str2, "moto")
}

func TestApplyColumnSorting(t *testing.T) {
	str1, str2 := applyColumnSorting("first second third", "uno dos tres", 2)
	assert.Equal(t, str1, "second")
	assert.Equal(t, str2, "dos")
	str1, str2 = applyColumnSorting("first second third", "uno dos tres", 3)
	assert.Equal(t, str1, "third")
	assert.Equal(t, str2, "tres")
	str1, str2 = applyColumnSorting("1 2 3 4 5", "5 4 3 2 1", 4)
	assert.Equal(t, str1, "4")
	assert.Equal(t, str2, "2")
}

func TestMySort(t *testing.T) {
	var testSlice = make([]string, 0)
	testSlice = append(testSlice, "1 4 3")
	testSlice = append(testSlice, "8 9 5")
	testSlice = append(testSlice, "6 2 7")
	var expectedSlice = make([]string, 0)
	expectedSlice = append(expectedSlice, "8 9 5")
	expectedSlice = append(expectedSlice, "6 2 7")
	expectedSlice = append(expectedSlice, "1 4 3")
	_ = mySort(&testSlice, true, true, true, 1)
	assert.Equal(t, expectedSlice, testSlice)

	testSlice = make([]string, 0)
	testSlice = append(testSlice, "1 4 3")
	testSlice = append(testSlice, "8 9 5")
	testSlice = append(testSlice, "6 2 7")
	expectedSlice = make([]string, 0)
	expectedSlice = append(expectedSlice, "6 2 7")
	expectedSlice = append(expectedSlice, "1 4 3")
	expectedSlice = append(expectedSlice, "8 9 5")
	_ = mySort(&testSlice, true, false, true, 2)
	assert.Equal(t, expectedSlice, testSlice)

	testSlice = make([]string, 0)
	testSlice = append(testSlice, "11 42 33")
	testSlice = append(testSlice, "84 95 56")
	testSlice = append(testSlice, "67 28 79")
	expectedSlice = make([]string, 0)
	expectedSlice = append(expectedSlice, "84 95 56")
	expectedSlice = append(expectedSlice, "67 28 79")
	expectedSlice = append(expectedSlice, "11 42 33")
	_ = mySort(&testSlice, true, true, false, 0)
	assert.Equal(t, expectedSlice, testSlice)

	testSlice = make([]string, 0)
	testSlice = append(testSlice, "Napkin")
	testSlice = append(testSlice, "Apple")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "BOOK")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "Hauptbahnhof")
	testSlice = append(testSlice, "Book")
	testSlice = append(testSlice, "Go")
	expectedSlice = make([]string, 0)
	expectedSlice = append(expectedSlice, "Napkin")
	expectedSlice = append(expectedSlice, "January")
	expectedSlice = append(expectedSlice, "January")
	expectedSlice = append(expectedSlice, "Hauptbahnhof")
	expectedSlice = append(expectedSlice, "Go")
	expectedSlice = append(expectedSlice, "Book")
	expectedSlice = append(expectedSlice, "BOOK")
	expectedSlice = append(expectedSlice, "Apple")
	_ = mySort(&testSlice, false, true, false, 0)
	assert.Equal(t, expectedSlice, testSlice)

	testSlice = make([]string, 0)
	testSlice = append(testSlice, "Napkin 2")
	testSlice = append(testSlice, "Apple 1")
	testSlice = append(testSlice, "January 3")
	testSlice = append(testSlice, "BOOK 4")
	testSlice = append(testSlice, "January 5")
	testSlice = append(testSlice, "Hauptbahnhof 6")
	testSlice = append(testSlice, "Book 7")
	testSlice = append(testSlice, "Go 0")
	expectedSlice = make([]string, 0)
	expectedSlice = append(expectedSlice, "Book 7")
	expectedSlice = append(expectedSlice, "Hauptbahnhof 6")
	expectedSlice = append(expectedSlice, "January 5")
	expectedSlice = append(expectedSlice, "BOOK 4")
	expectedSlice = append(expectedSlice, "January 3")
	expectedSlice = append(expectedSlice, "Napkin 2")
	expectedSlice = append(expectedSlice, "Apple 1")
	expectedSlice = append(expectedSlice, "Go 0")
	err := mySort(&testSlice, false, true, true, 2)
	assert.Equal(t, expectedSlice, testSlice)
	assert.Nil(t, err, "error in sort non-numeric strings")

}

func TestValidate(t *testing.T) {
	testSlice := make([]string, 0)
	testSlice = append(testSlice, "1 4")
	testSlice = append(testSlice, "2 5")
	testSlice = append(testSlice, "3 6")
	err := validate(5, 1, true, testSlice)
	assert.NotNil(t, err, "Columns error")
	testSlice = make([]string, 0)
	testSlice = append(testSlice, "Napkin")
	testSlice = append(testSlice, "Apple")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "BOOK")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "Hauptbahnhof")
	testSlice = append(testSlice, "Book")
	testSlice = append(testSlice, "Go")
	assert.NotNil(t, validate(0, 1, true, testSlice), "numbersFile error")
}

func printLinesToFile(filePath string, values []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range values {
		_, _ = fmt.Fprintln(f, value)
	}
	return nil
}

func TestFileInputAndOutput(t *testing.T) {
	testSlice := make([]string, 0)
	testSlice = append(testSlice, "Napkin")
	testSlice = append(testSlice, "Apple")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "BOOK")
	testSlice = append(testSlice, "January")
	testSlice = append(testSlice, "Hauptbahnhof")
	testSlice = append(testSlice, "Book")
	testSlice = append(testSlice, "Go")
	_ = output(testSlice, "testFileOutput.txt", false, false)
	file, err := os.Open("testFile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	stringsFromFile, _, err := readStrings(file)
	assert.Equal(t, testSlice, stringsFromFile)
}
