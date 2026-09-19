package main

import (
	"log"
	"os"
)

func main() {
	logger := getLogger("/tmp/learning-lsp.log")

	logger.Println("hello world")
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to open learning-lsp log file")
	}

	return log.New(logfile, "[learning-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
