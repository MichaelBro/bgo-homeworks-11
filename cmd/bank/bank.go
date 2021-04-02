package main

import (
	"bgo-homeworks-11/cmd/bank/app"
	"bgo-homeworks-11/pkg/card"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	cardSvc := card.NewService()
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := &http.Server{
		Addr: addr,
		Handler: application,
	}
	return server.ListenAndServe()
}