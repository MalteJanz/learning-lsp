package main

import (
	"bufio"
	"encoding/json"
	"io"
	"learning-lsp/lsp"
	"learning-lsp/rpc"
	"log"
	"os"
)

func main() {
	logger := getLogger("/tmp/learning-lsp.log")
	logger.Println("lsp server is booting up...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	for scanner.Scan() {
		rawMsg := scanner.Bytes()
		method, contents, err := rpc.DecodeMsg(rawMsg)
		if err != nil {
			logger.Printf("Got an error while decoding message: %s. Skipping the message", err)
			continue
		}

		handleMsg(logger, writer, method, contents)
	}

	if err := scanner.Err(); err != nil {
		logger.Printf("Error reading input: %s", err)
	}
}

func handleMsg(logger *log.Logger, writer io.Writer, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse intitialize: %s", err)
		}

		logger.Printf("LSP client %s %s connected", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitializeResponse(request.ID)

		SendMsg(logger, writer, msg)
	case "initialized":
		logger.Println("Initialization completed")
	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Failed to parse textDocument/didOpen: %s", err)
		}

		logger.Printf("Opened: %s", notification.Params.TextDocument.URI)
	}
}

func SendMsg(logger *log.Logger, writer io.Writer, msg any) {
	data := rpc.EncodeMsg(msg)
	n, err := writer.Write([]byte(data))
	if err != nil {
		logger.Printf("Failed to write msg to stdout: %s", err)
		return
	}

	if len(data) != n {
		logger.Printf("Only partially written msg to stdout, %d out of %d bytes", n, len(data))
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to open learning-lsp log file")
	}

	return log.New(logfile, "[learning-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
