package domains

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TextResult wraps a plain text string into the CallToolResult structure
// expected by MCP tool handlers.
func TextResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, nil, nil
}

// ResourceResult constructs a ReadResourceResult with a single text content
// entry. Use this in MCP resource template handlers.
func ResourceResult(uri, mimeType, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			URI:      uri,
			MIMEType: mimeType,
			Text:     text,
		}},
	}
}

// PromptResult constructs a GetPromptResult from a description and one or more
// prompt messages. Use this in MCP prompt handlers.
func PromptResult(description string, messages ...*mcp.PromptMessage) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: description,
		Messages:    messages,
	}
}

// UserMessage creates a user-role PromptMessage with plain text content.
func UserMessage(text string) *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Role:    "user",
		Content: &mcp.TextContent{Text: text},
	}
}
