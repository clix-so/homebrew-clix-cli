package logx

import (
	"fmt"
	"strings"
	"time"
)

const (
	Bold      = "1"
	Gray      = "90"
	Reset 	  = "\033[0m"
)

type Logger struct {	
	indent int
	prefix string
	codes  []string
	useSpinner bool
}

func Log() *Logger {
	return &Logger{}
}

func (l *Logger) Indent(n int) *Logger {
	l.indent = n
	return l
}

func (l *Logger) Branch() *Logger {
	l.prefix = l.prefix + " └ "
	return l
}

func (l *Logger) Gray() *Logger {
	l.codes = append(l.codes, Gray)
	return l
}

func (l *Logger) Bold() *Logger {
	l.codes = append(l.codes, Bold)
	return l
}

func (l *Logger) Code() *Logger {
	l.Gray()
	return l
}

func (l *Logger) Title() *Logger {
	l.Gray().Bold()
	return l
}

func (l *Logger) Success() *Logger {
	l.prefix = l.prefix + "✅ "
	return l
}

func (l *Logger) Failure() *Logger {
	l.prefix = l.prefix + "❌ "
	return l
}

func (l *Logger) WithSpinner() *Logger {
	l.useSpinner = true
	return l
}

func (l *Logger) Println(msg string) {
	style := ansiCode(l.codes...)
	indent := spaces(l.indent)

	lines := strings.Split(msg, "\n")
	var fullLines []string

	for i, line := range lines {
		if i == 0 {
			fullLines = append(fullLines, fmt.Sprintf("%s%s%s%s", indent+l.prefix, style, line, Reset))
		} else {
			fullLines = append(fullLines, fmt.Sprintf("%s%s%s%s", indent, style, line, Reset))
		}
	}

	full := strings.Join(fullLines, "\n")

	if l.useSpinner {
		spinner := NewSpinner()
		spinner.Start(full)
		time.Sleep(1 * time.Second)
		spinner.Stop()
	} else {
		fmt.Println(full)
	}
}

func Separatorln() {
	fmt.Println("─────────────────────────────────────────────────────")
}

func NewLine() {
	fmt.Println()
}

func ansiCode(codes ...string) string {
	if len(codes) == 0 {
		return ""
	}
	return "\033[" + strings.Join(codes, ";") + "m"
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}
