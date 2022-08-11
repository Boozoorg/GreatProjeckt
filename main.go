package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Boozoorg/GreatProjeck/accounts"
	"github.com/Boozoorg/GreatProjeck/app"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/dig"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://postgres:postgres@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(0)
	}
}

func execute(host, port, dsn string) (err error) {
	deps := []interface{}{
		app.NewServeMux,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		accounts.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	conteiner := dig.New()
	for _, dep := range deps {
		err := conteiner.Provide(dep)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = conteiner.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		log.Println(err)
		return err
	}

	err = conteiner.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
	return
}
