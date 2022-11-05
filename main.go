package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var message string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input:")
	scanner.Scan()
	input := scanner.Text()
	for i := 0; i < len(input); i++ {
		message += scrambler(string(input[i]))
	}
	fmt.Println(message)
}

func scrambler(inputedLetter string) string {
	scrambles := map[interface{}]interface{}{
		"h": "q",
		"i": "p",
	}
	if scrambles[inputedLetter] != nil {
		return scrambles[inputedLetter].(string)
	} else {
		return inputedLetter
	}
}
