package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Boozoorg/GreatProjeck/accounts"
	"github.com/Boozoorg/GreatProjeck/app"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		log.Print(err)
		os.Exit(0)
	}
}

func execute(host, port string) (err error) {
	mux := http.NewServeMux()
	AccountService := accounts.NewService()
	server := app.NewServeMux(mux, AccountService)
	server.Init()

	srv := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}
