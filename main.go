package main

import (
	"start_point/accounts"
	"start_point/app"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		log.Print(err)
		os.Exit(2)
	}
}

func execute(host, port string) (err error) {
	mux  := http.NewServeMux()
	AccountService := accounts.NewService()
	server := app.NewServeMux(mux, AccountService)
	server.Init()
	
	srv := http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}