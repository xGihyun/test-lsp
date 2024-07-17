package lsp

type InitRequest struct {
	Request
	Params InitRequestParams `json:"params"`
}

type InitRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`

	// A bunch more stuff here
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitResponse struct {
	Response
	Result InitResult `json:"result"`
}

type InitResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync   int            `json:"textDocumentSync"`
	HoverProvider      bool           `json:"hoverProvider"`
	DefinitionProvider bool           `json:"definitionProvider"`
	CodeActionProvider bool           `json:"codeActionProvider"`
	CompletionProvider map[string]any `json:"completionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitResponse(id int) InitResponse {

	return InitResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true, // Not implemented
				CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "testlsp",
				Version: "0.0.1-beta",
			},
		},
	}
}
