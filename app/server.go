package app

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/Boozoorg/GreatProjeck/accounts"
)

type Server struct {
	mux *http.ServeMux
	accountService *accounts.Service
}

func NewServeMux(Mux *http.ServeMux, AccountService *accounts.Service) *Server {
	return &Server{
		mux: Mux,
		accountService: AccountService,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s Server) Init() {
	s.mux.HandleFunc("/account.save_acc", s.SaveAccount)
	s.mux.HandleFunc("/account.send_message", s.SendMessage)
}

func (s *Server) SaveAccount(writer http.ResponseWriter, request *http.Request) {
	Name := request.URL.Query().Get("name")
	Password := request.URL.Query().Get("password")

	item, err := s.accountService.Registration(Name, Password)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	data, err := json.Marshal(item)
	log.Print(string(data))
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) SendMessage(writer http.ResponseWriter, request *http.Request) {
	From := request.URL.Query().Get("from")
	To := request.URL.Query().Get("to")
	Message := request.URL.Query().Get("message")

	item, err := s.accountService.SendMessage(From, To, Message)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	data, err := json.Marshal(item)
	log.Print(string(data))
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}