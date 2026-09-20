package analysis

import (
	"fmt"
	"learning-lsp/lsp"
	"strings"
	"unicode/utf8"
)

func (s *State) Hover(req lsp.HoverRequest) *lsp.HoverResult {
	uri := req.Params.TextDocument.URI
	if !strings.HasSuffix(uri, ".twig") && !strings.HasSuffix(uri, ".js") && !strings.HasSuffix(uri, ".ts") {
		return nil
	}

	line := req.Params.Position.Line
	char := req.Params.Position.Character
	fileContent := s.Documents[uri]

	index, err := positionToOffset(fileContent, line, char)
	if err != nil {
		s.Logger.Printf("Failed to convert position to offset: %s", err)
		return nil
	}

	quotedString := expandToQuoteOrLineBound(fileContent, index)
	if quotedString == "" {
		return nil
	}

	if snippets, ok := s.Snippets[quotedString]; ok {
		return &lsp.HoverResult{
			Contents: fmt.Sprintf("defined as: %v", snippets),
		}
	}

	return nil
}

// Disclaimer: Build by Claude
// positionToOffset converts an LSP line/character position (UTF-16 code units)
// into a byte offset into content (UTF-8).
func positionToOffset(content string, line, char int) (int, error) {
	// 1. Find the byte offset where the target line starts.
	lineStart := 0
	for l := 0; l < line; l++ {
		idx := strings.IndexByte(content[lineStart:], '\n')
		if idx == -1 {
			return 0, fmt.Errorf("line %d out of range", line)
		}
		lineStart += idx + 1
	}

	// 2. Walk runes on that line, counting UTF-16 units, until we've
	//    consumed `char` of them.
	rest := content[lineStart:]
	units := 0
	byteOffset := lineStart

	for units < char {
		r, size := utf8.DecodeRuneInString(rest)
		if r == utf8.RuneError && size == 0 {
			// ran off the end of the string (position past EOF)
			return byteOffset, nil
		}
		if r == '\n' {
			// position is past the end of this line; clamp here
			return byteOffset, nil
		}

		if r > 0xFFFF {
			units += 2 // encoded as a UTF-16 surrogate pair
		} else {
			units++
		}
		byteOffset += size
		rest = rest[size:]
	}

	return byteOffset, nil
}

func expandToQuoteOrLineBound(content string, offset int) string {
	leftPos, leftFound := scanLeft(content, offset)
	rightPos, rightFound := scanRight(content, offset)

	if leftFound && rightFound {
		return content[leftPos:rightPos]
	}

	return ""
}

// returned pos is inclusive (pointing right after the quote if found)
func scanLeft(content string, offset int) (pos int, found bool) {
	for i := offset; i > 0; i-- {
		c := content[i-1]
		if c == '\n' {
			return i, false
		}
		if c == '\'' || c == '"' {
			return i, true
		}
	}

	return 0, false
}

// returned pos is exclusive (pointing at the quote if found)
func scanRight(content string, offset int) (pos int, found bool) {
	for i := offset; i < len(content); i++ {
		c := content[i]
		if c == '\n' {
			return i, false
		}
		if c == '\'' || c == '"' {
			return i, true
		}
	}

	return 0, false
}
