package analysis

import (
	"encoding/json"
	"strings"
)

func (s *State) indexFile(uri, text string) {
	if !strings.HasSuffix(uri, "snippet/en.json") && !strings.HasSuffix(uri, "snippet/de.json") {
		return
	}

	s.Logger.Printf("indexing file %s", uri)

	// Todo: build a lexer + parser from scratch to also keep line + column info as that is important for definition lookup
	var root interface{}
	if err := json.Unmarshal([]byte(text), &root); err != nil {
		s.Logger.Printf("json unmarshal failed for %s with err: %s", uri, err)
		return
	}

	flattenInsert(uri, root, s.Snippets)
}

type stackItem struct {
	node interface{}
	path string
}

func flattenInsert(uri string, root interface{}, snippets map[string][]SnippetValue) {
	stack := []stackItem{
		{node: root, path: ""},
	}

	for len(stack) > 0 {
		// pop
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch v := item.node.(type) {
		case map[string]interface{}:
			for key, val := range v {
				childPath := key
				if item.path != "" {
					childPath = item.path + "." + key
				}
				stack = append(stack, stackItem{node: val, path: childPath})
			}
		case string:
			// leaf, only considering strings
			snippets[item.path] = append(snippets[item.path], SnippetValue{
				URI:   uri,
				Value: v,
			})
		}
	}
}
