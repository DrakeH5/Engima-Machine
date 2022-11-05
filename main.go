package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input:")
	scanner.Scan()
	input := scanner.Text()
	var message string = input
	fmt.Println(message)

}
