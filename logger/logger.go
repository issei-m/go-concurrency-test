package logger

import (
	"fmt"
	"os"
)

func Debug(message string) {
	fmt.Println(fmt.Sprintf("\033[37m%s\033[0m", message))
}

func Info(message string) {
	fmt.Println(message)
}

func Error(message string) {
	if _, err := fmt.Fprintln(os.Stderr, fmt.Sprintf("\033[31m%s\033[0m", message)); err != nil {
		panic(err)
	}
}
