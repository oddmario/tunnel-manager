package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func Init() {
	logger = log.New(os.Stdout, "[LOG] ", log.LstdFlags)
}

func Warn(v ...any) {
	args := append([]any{"[WARN]"}, v...)

	logger.Println(args...)
}

func Info(v ...any) {
	args := append([]any{"[INFO]"}, v...)

	logger.Println(args...)
}

func Error(v ...any) {
	args := append([]any{"[ERROR]"}, v...)

	logger.Println(args...)
}

func Fatal(v ...any) {
	args := append([]any{"[FATAL]"}, v...)

	logger.Fatalln(args...)
}
