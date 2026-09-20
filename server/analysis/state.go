package analysis

import "log"

type State struct {
	Logger   *log.Logger
	RootPath string
	// Map of file names to contents
	Documents map[string]string
	// Map of snippet key to array of SnippetValue (e.g. can be defined in multiple languages / files)
	Snippets map[string][]SnippetValue
}

type SnippetValue struct {
	// File URI where the snippet is defined
	URI string
	// translated value
	Value string
}

func NewState(logger *log.Logger) State {
	return State{
		Logger:    logger,
		RootPath:  "",
		Documents: map[string]string{},
		Snippets:  map[string][]SnippetValue{},
	}
}

func (s *State) Initialize(rootPath string) {
	s.RootPath = rootPath
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
	s.indexFile(uri, text)
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
	s.indexFile(uri, text)
}
