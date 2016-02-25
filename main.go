package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	flagListenPort  = 8080
	flagListenIface = "0.0.0.0"
	flagListenAddr  = ""
)

func main() {
	serve := flag.Bool("s", false, "serve sausages")
	flag.StringVar(&flagListenIface, "i", flagListenIface, "list on interface")
	flag.IntVar(&flagListenPort, "p", flagListenPort, "list on TCP port")
	flag.Parse()

	flagListenAddr = fmt.Sprintf("%s:%d", flagListenIface, flagListenPort)

	http.HandleFunc("/", serveSausages)

	if *serve {
		panicon(http.ListenAndServe(flagListenAddr, nil))
	} else {
		flag.PrintDefaults()
	}
}

func panicon(err error) {
	if err != nil {
		panic(err)
	}
}
