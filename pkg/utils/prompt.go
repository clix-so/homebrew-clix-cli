package utils

import (
	"fmt"
)

func Prompt(message string) string {
	fmt.Print(message + ": ")
	var input string
	fmt.Scanln(&input)
	return input
}
