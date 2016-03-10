package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	flagListenPort  = 8080
	flagListenIface = "0.0.0.0"
	flagListenAddr  = ""

	flagClientWorkers = 8192
)

var (
	serve = kingpin.Command("serve", "serve sausages generously")

	sizzle          = kingpin.Command("sizzle", "consume sausages greedily")
	sizzleProxyURL  = sizzle.Arg("proxy", "proxy service URL").Default("http://localhost:8080/").URL()
	sizzleWorkers   = sizzle.Flag("workers", "worker count").Default("8192").Int()
	sizzleStaggerBy = sizzle.Flag("stagger", "stagger the start of each worker by ms").Default("10").Int()
)

func main() {
	switch kingpin.Parse() {
	case "serve":
		ActionServe()

	case "sizzle":
		ActionSizzle()

	default:
		panicon(fmt.Errorf("Nope"))
	}
}

func usage(code int) {

	os.Exit(code)
}

func ActionServe() {
	srv := NewSausageServer()
	panicon(srv.ListenAndServe())
}

func ActionSizzle() {
	fmt.Printf("Sizzling %v\n", *sizzleProxyURL)

	fields := []string{"", "", "bytes", "", "", "duration", "method", "url", "", "", "mime_type", "agent"}
	l := NewSSVLexer(fields)
	src := NewLogParser(os.Stdin, l)

	client := NewClient(*sizzleProxyURL, src)
	client.Workers = *sizzleWorkers
	client.StaggerBy = time.Millisecond * time.Duration(*sizzleStaggerBy)
	client.Sizzle()
}

func panicon(err error) {
	if err != nil {
		panic(err)
	}
}
