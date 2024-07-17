package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"test-lsp/analysis"
	"test-lsp/lsp"
	"test-lsp/rpc"
)

func main() {
	logger := getLogger("/home/gihyun/Documents/Programming/Go/test-lsp/log.txt")
	logger.Println("LSP Starting...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)

		if err != nil {
			logger.Printf("Error: %s", err)
			continue
		}

		handleMessage(logger, writer, &state, method, content)
	}
}

func writeResponse(writer io.Writer, message any) {
	reply := rpc.EncodeMessage(message)
	writer.Write([]byte(reply))
}

func handleMessage(logger *log.Logger, writer io.Writer, state *analysis.State, method string, contents []byte) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitRequest

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		msg := lsp.NewInitResponse(request.ID)
		writeResponse(writer, msg)

		logger.Print("Reply sent.")
		break

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		notification := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{RPC: "2.0", Method: "textDocument/publishDiagnostics"},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		}

		writeResponse(writer, notification)
		break

	case "textDocument/didChange":
		var request lsp.DidChangeNotification

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			notification := lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{RPC: "2.0", Method: "textDocument/publishDiagnostics"},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			}

			writeResponse(writer, notification)
		}
		break

	case "textDocument/hover":
		var request lsp.HoverRequest

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
		break

	case "textDocument/definition":
		var request lsp.DefinitionRequest

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
		break

	case "textDocument/codeAction":
		// TODO: Stuff
		break

	case "textDocument/completion":
		var request lsp.CompletionRequest

		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed to parse contents: %s", err)
		}

		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)
		break
	}

}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("Logger file not found.")
	}

	return log.New(logfile, "[testinglsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
