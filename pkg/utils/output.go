package utils

import (
	"fmt"
	"time"
)

const (
	Bold      = "\033[1m"
	Gray      = "\033[90m"
	Reset     = "\033[0m"
)

func Separatorln() {
	fmt.Println("─────────────────────────────────────────────────────")
}

func Grayln(msg string) {
	fmt.Printf("\033[90m%s\033[0m\n", msg)
}

func Boldln(msg string) {
	fmt.Printf("\033[1m%s\033[0m\n", msg)
}

func GrayBoldln(msg string) {
	fmt.Printf("\033[90;1m%s\033[0m\n", msg)
}

func Titleln(msg string) {
	GrayBoldln(msg)
}

func Code(msg string) {
	Grayln(msg)
}

func Successln(msg string) {
	fmt.Println("✅ " + msg)
}

func Failureln(msg string) {
	fmt.Println("❌ " + msg)
}

func Warnln(msg string) {
	fmt.Println("⚠️  " + msg)
}

func TitlelnWithSpinner(msg string) {
	spinner := NewSpinner()
	spinner.Start(msg)
	time.Sleep(1 * time.Second)
	spinner.Stop()
}

func Indentln(msg string, spaces int) {
	for i := 0; i < spaces; i++ {
		fmt.Print(" ")
	}
	fmt.Println(msg)
}

func Branchln(msg string) {
	fmt.Println(" └ " + msg)
}

func BranchSuccessln(msg string) {
	fmt.Println(" └ ✅ " + msg)
}

func BranchFailureln(msg string) {
	fmt.Println(" └ ❌ " + msg)
}