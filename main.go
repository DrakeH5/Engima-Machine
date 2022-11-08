package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var rotors [3]map[interface{}]interface{} = generateRotors()

func main() {
	var message string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input:")
	scanner.Scan()
	input := scanner.Text()
	for i := 0; i < len(input); i++ {
		message += scrambler(string(input[i]))
	}
	message = goThroughRotors(message)
	var output string
	for i := 0; i < len(message); i++ {
		output += scrambler(string(message[i]))
	}
	fmt.Println(output)
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

func generateRotors() [3]map[interface{}]interface{} {
	rand.Seed(time.Now().UnixNano())
	rotors := [3]map[interface{}]interface{}{}
	letters := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for j := 0; j < 3; j++ {
		unusedLetters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
		rotor := map[interface{}]interface{}{
			"a": "",
			"b": " ",
			"c": " ",
			"d": " ",
			"e": " ",
			"f": " ",
			"g": " ",
			"h": " ",
			"i": " ",
			"j": " ",
			"k": " ",
			"l": " ",
			"m": " ",
			"n": " ",
			"o": " ",
			"p": " ",
			"q": " ",
			"r": " ",
			"s": " ",
			"t": " ",
			"u": " ",
			"v": " ",
			"w": " ",
			"x": " ",
			"y": " ",
			"z": " ",
		}
		for i := 0; i < 26; i++ {
			var randNmb int = rand.Intn(len(unusedLetters)-0) + 0
			rotor[letters[i]] = unusedLetters[randNmb]
			unusedLetters = RemoveIndex(unusedLetters[:], randNmb)
		}
		rotors[j] = rotor
	}
	return rotors
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func generateReflector() map[interface{}]interface{} {
	letters := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	unusedLetters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	reflector := map[interface{}]interface{}{
		"a": "",
		"b": " ",
		"c": " ",
		"d": " ",
		"e": " ",
		"f": " ",
		"g": " ",
		"h": " ",
		"i": " ",
		"j": " ",
		"k": " ",
		"l": " ",
		"m": " ",
		"n": " ",
		"o": " ",
		"p": " ",
		"q": " ",
		"r": " ",
		"s": " ",
		"t": " ",
		"u": " ",
		"v": " ",
		"w": " ",
		"x": " ",
		"y": " ",
		"z": " ",
	}
	for i := 0; i < 26; i++ {
		var randNmb int = rand.Intn(len(unusedLetters))
		reflector[letters[i]] = unusedLetters[randNmb]
		unusedLetters = RemoveIndex(unusedLetters[:], randNmb)
	}
	return reflector
}

func goThroughRotors(input string) string {
	var output string
	var reflector map[interface{}]interface{} = generateReflector()
	for i := 0; i < len(input); i++ {
		if string(input[i]) != " " {
			output += rotorConversions(string(input[i]), reflector)
		} else {
			output += " "
		}
	}
	return output
}

func rotorConversions(input string, reflector map[interface{}]interface{}) string {
	var output string = input
	for i := 0; i < 3; i++ {
		output = rotors[i][output].(string)
		rotateRotors()
		if i == 2 {
			output = reflector[output].(string)
			rotateRotors()
			for i := 2; i > -1; i-- {
				output = rotors[i][output].(string)
				rotateRotors()
			}
		}
	}
	return output
}

var nbmOfRotations int = 0

func rotateRotors() {
	letters := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	last := rotors[0]["z"]
	for i := 25; i > 0; i-- {
		rotors[0][letters[i]] = rotors[0][letters[i-1]]
	}
	rotors[0]["a"] = last
	if nbmOfRotations%26 == 0 {
		last1 := rotors[1]["z"]
		for i := 25; i > 0; i-- {
			rotors[1][letters[i]] = rotors[1][letters[i-1]]
		}
		rotors[1]["Z"] = last1
	}
	if nbmOfRotations%52 == 0 {
		last2 := rotors[2]["z"]
		for i := 25; i > 0; i-- {
			rotors[2][letters[i]] = rotors[2][letters[i-1]]
		}
		rotors[2]["a"] = last2
	}
	nbmOfRotations++
}
