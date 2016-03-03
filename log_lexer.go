package main

// LogLexer is an interface which reads an access log line from an io.Reader
// and tokenizes it into a key/value map.
type LogLexer interface {
	Lex(string) (map[string]string, error)
}
