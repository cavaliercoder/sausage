package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	flagListenPort  = 8080
	flagListenIface = "0.0.0.0"
	flagListenAddr  = ""
)

func main() {
	if len(os.Args) < 2 {
		panicon(fmt.Errorf("more!"))
	}

	switch os.Args[1] {
	case "serve":
		serve()

	case "sizzle":
		proxyURL := "http://localhost:8080/"
		if len(os.Args) > 2 {
			proxyURL = os.Args[2]
		}

		sizzle(proxyURL)

	default:
		panicon(fmt.Errorf("Nope"))
	}
}

func serve() {
	flagListenAddr = fmt.Sprintf("%s:%d", flagListenIface, flagListenPort)
	http.HandleFunc("/", serveSausages)
	panicon(http.ListenAndServe(flagListenAddr, nil))
}

func sizzle(proxyURL string) {
	fmt.Printf("Sizzling %s\n", proxyURL)

	u, err := url.Parse(proxyURL)
	panicon(err)

	fields := []string{"", "", "bytes", "", "", "duration", "method", "url", "", "", "mime_type", "agent"}
	l := NewSSVLexer(fields)
	src := NewLogParser(os.Stdin, l)

	client := NewClient(u, src)
	client.Workers = 1000
	client.Sizzle()
}

func panicon(err error) {
	if err != nil {
		panic(err)
	}
}
