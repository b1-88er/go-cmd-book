package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	host := flag.String("h", "localhost", "host")
	port := flag.Int("p", 8080, "port")
	todoFile := flag.String("f", "todoServer.json", "todo JSON file")

	flag.Parse()

	s := http.Server{
		Addr:         fmt.Sprintf("%s:%d", *host, *port),
		Handler:      newMux(*todoFile),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
