package analysis

import (
	"fmt"
	"strings"
	"test-lsp/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnostics(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnostics(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// This should look up the type

	documents := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(documents)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				// NOTE: This indicates where the definition is, hard coded for educational purposes
				Start: lsp.Position{
					Line:      position.Line + 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line + 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Honkai Impact 3rd",
			Detail:        "Tuna",
			Documentation: "(´｡• ᵕ •｡`)",
		},
	}

	// Ask static analysis tools to figure out good completions
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
}

func getDiagnostics(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "Error") {
			idx := strings.Index(line, "Error")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(uint(row), uint(idx), uint(idx+len("Error"))),
				Severity: 1,
				Source:   "Imagination",
				Message:  "Fix this error! (๑˃ᴗ˂)ﻭ",
			})
		}
	}

	return diagnostics
}

func LineRange(line, start, end uint) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
