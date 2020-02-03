package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ernestosuarez/itertools"
)

var iterable = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var baseSlice = make([][]int, 0)
var loopCounter int
var exit bool

type possibleGuesses struct {
	possibleValue []int
	count         int
}

func setNummberOfDigits() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Set the number of digits")
	input, _ := reader.ReadString('\n')
	digitNumber := strings.TrimSuffix(input, "\n")
	maxLetters, err := strconv.Atoi(digitNumber)
	if err != nil {
		fmt.Println("cannot convert to number")
		os.Exit(2)
	}
	return maxLetters
}

func setMultiSlice(s int) [][]int {

	for v := range itertools.PermutationsInt(iterable, s) {
		baseSlice = append(baseSlice, v)
	}
	return baseSlice
}

func statistics(currentArray [][]int) {
	fmt.Println(currentArray)
	fmt.Println("Possible combinations:")
	fmt.Print(">>> ")
	fmt.Println(len(currentArray))
}

func askGuess(counter int) {
	fmt.Print()
	if counter < 2 {
		fmt.Print("Your first guess will be?")
	} else {
		fmt.Printf("Guess n. %v will be?", counter)
	}
	fmt.Println()
}

func getGuess() []int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>>")
	input, _ := reader.ReadString('\n')
	digitNumber := strings.TrimSuffix(input, "\n")
	newSlice := make([]int, 0)
	for _, v := range digitNumber {
		number, _ := strconv.Atoi(string(v))
		newSlice = append(newSlice, number)
	}
	fmt.Println(newSlice)
	return newSlice
}

func getBullsCows(digitNumber int) string {
	fmt.Println("Input BullsCow combination like bcc")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	cowsbulls := strings.TrimSuffix(input, "\n")
	if len(cowsbulls) > digitNumber {
		fmt.Println("wrong input!")
		getBullsCows(digitNumber)
	}
	return cowsbulls
}

func isAllBulls(string string, digitNumber int) bool {
	for _, v := range string {
		if v != 'b' {
			return false
		}
	}
	return true
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func filterSlice(guess []int, s string, currentArray [][]int, digitNumber int) ([][]int, bool) {
	if (len(s) == digitNumber) && isAllBulls(s, digitNumber) {
		fmt.Println("Game over, you won!")
		return currentArray, true
	}
	bullsCowsSlice := strings.Fields(s)
	newCurrentArray := make([][]int, 0)
	for _, v := range currentArray {
		sliceValue := make([]string, 0)
		for i, c := range v {
			if contains(guess, c) && guess[i] == c {
				sliceValue = append(sliceValue, "b")
			} else if contains(guess, c) {
				sliceValue = append(sliceValue, "c")
			}
		}

		if strings.Join(bullsCowsSlice, "") == strings.Join(sliceValue, "") {
			newCurrentArray = append(newCurrentArray, v)
		}
	}
	return newCurrentArray, false
}

func allPossibilities(digitNumber int) []string {
	possibilites := make([]string, 0)
	for combination := range generateCombinations("bc", digitNumber) {
		combination = sortString(combination)
		possibilites = append(possibilites, combination)
		possibilites = sliceUniqMap(possibilites)
	}
	return possibilites
}

func sortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func sliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)
	go func(c chan string) {
		defer close(c)
		addLetter(c, "", alphabet, length)
	}(c)
	return c
}

func addLetter(c chan string, combo string, alphabet string, length int) {
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

func calculateNextGuess(currentArray [][]int, baseArray [][]int, possibilities []string) []int {

	possibleGuessesCollection := []possibleGuesses{}
	for _, p := range currentArray {
		pGuess := possibleGuesses{p, 0}
		count := 0
		for _, e := range baseArray {
			for _, i := range e {
				if contains(p, i) {
					count++
					pGuess = possibleGuesses{p, count}
				}
			}
		}
		possibleGuessesCollection = append(possibleGuessesCollection, pGuess)
	}
	vv := minIntSlice(possibleGuessesCollection)
	return vv.possibleValue
}

func minIntSlice(v []possibleGuesses) possibleGuesses {
	m := 10000
	ff := possibleGuesses{[]int{0}, 0}
	for _, e := range v {
		if e.count < m {
			ff = e
		}
	}
	fmt.Println(ff)
	return ff
}

func main() {
	digitNumber := setNummberOfDigits()
	baseArray := setMultiSlice(digitNumber)
	currentArray := baseArray
	possibilities := allPossibilities(digitNumber)
outer:
	for {
		loopCounter++
		statistics(currentArray)
		askGuess(loopCounter)
		newGuess := getGuess()
		bullscows := getBullsCows(digitNumber)
		currentArray, exit = filterSlice(newGuess, bullscows, currentArray, digitNumber)
		if exit == true {
			break outer
		}
		calculateNextGuess(currentArray, baseArray, possibilities)
	}
}
