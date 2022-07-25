package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	s.mux.HandleFunc("/account.save", s.SaveAccount)
	s.mux.HandleFunc("/account.delate", s.DelateAccByID)
}

func (s *Server) SaveAccount(writer http.ResponseWriter, request *http.Request) {
	Name := request.URL.Query().Get("name")
	Password := request.URL.Query().Get("password")
	Mail := request.URL.Query().Get("mail")

	item, err := s.accountService.Registration(Name, Password, Mail)
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

func (s *Server) DelateAccByID(writer http.ResponseWriter, request *http.Request) {
	ID := request.URL.Query().Get("id")

	//Конвертируем string в uint64
	
	id, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}


	item, err := s.accountService.DelateAccByID(id)
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