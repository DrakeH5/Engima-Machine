package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	var rotors [3]map[interface{}]interface{} = generateRotors()
	var message string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input:")
	scanner.Scan()
	input := scanner.Text()
	for i := 0; i < len(input); i++ {
		message += scrambler(string(input[i]))
	}
	message = goThroughRotors(rotors, message)
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

func generateRotors() [3]map[interface{}]interface{} {
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
			var randNmb int = rand.Intn(len(unusedLetters))
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

func goThroughRotors(rotors [3]map[interface{}]interface{}, input string) string {
	var output string
	var reflector map[interface{}]interface{} = generateReflector()
	for i := 0; i < len(input); i++ {
		output += rotorConversions(rotors, string(input[i]), reflector)
	}
	return output
}

func rotorConversions(rotors [3]map[interface{}]interface{}, input string, reflector map[interface{}]interface{}) string {
	var output string
	for i := 0; i < 3; i++ {
		output = rotors[i][input].(string)
		if i == 2 {
			output = reflector[input].(string)
			for i := 2; i > -1; i-- {
				output = rotors[i][input].(string)
			}
		}
	}
	return output
}
