package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	RootPath   string      `json:"rootPath"`
	// ... there's tons more that goes here
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`

	HoverProvider bool `json:"hoverProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		RPC: "2.0",
		ID:  &id,
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1, // Full: Documents are synced by always sending the full content of the document.
				HoverProvider:    true,
			},
			ServerInfo: ServerInfo{
				Name:    "learning-lsp",
				Version: "0.1.0",
			},
		},
	}
}
