package main

import (
	"fmt"
	"sort"
	"os"
	"bufio"
)

var words = make(map[string]int)

func addWord(word string) {
	count, exist := words[word]

	if exist {
		words[word] = count + 1
	} else {
		words[word] = 1
	}
}

func printArray() {
	var keys []string

	for k := range words {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s\t\t%d\n", k, words[k])
	}
}

func readFile(name string) {
	file, error := os.Open(name)

	if error != nil {
		fmt.Printf("%s\n", error)
		return
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		addWord(word)
	}
}

func main() {
	args := os.Args

	for i := 1; i < len(args); i++ {
		readFile(args[i])
	}

	printArray()
}