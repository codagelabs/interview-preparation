package main

import (
	"fmt"
	"strings"
)

func reverseWords(s string) string {

	wordsArray := strings.Fields(s)

	for i, j := 0, len(wordsArray)-1; i < j; i, j = i+1, j-1 {
		wordsArray[i], wordsArray[j] = wordsArray[j], wordsArray[i]
	}
	return strings.Join(wordsArray, " ")

}

func main() {
	fmt.Println(reverseWords("hello i am rahul"))
}
