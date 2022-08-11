package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Boozoorg/GreatProjeck/client"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Server struct {
	mux           *mux.Router
	clientService *client.Service
}

func NewServeMux(Mux *mux.Router, AccountService *client.Service) *Server {
	return &Server{
		mux:           Mux,
		clientService: AccountService,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

const (
	GET  = "GET"
	POST = "POST"
)

func (s *Server) Init() {
	s.mux.HandleFunc("/registration", s.Registration).Methods(POST)
	s.mux.HandleFunc("/delete/{id}", s.DeleteAccountById).Methods(POST)
	s.mux.HandleFunc("/chat", s.Messanger).Methods(POST)
	s.mux.HandleFunc("/chat", s.GetChatStory).Methods(GET)
}

var item *client.Account

func (s *Server) Registration(writer http.ResponseWriter, request *http.Request) {
	err := json.NewDecoder(request.Body).Decode(&item)

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err = s.clientService.Registration(request.Context(), item)
	if errors.Is(err, client.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if errors.Is(err, client.ErrAlreadyExe) {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
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

func (s *Server) DeleteAccountById(writer http.ResponseWriter, request *http.Request) {
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.clientService.DeleteAccount(request.Context(), item)
	if errors.Is(err, client.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Messanger(writer http.ResponseWriter, request *http.Request) {
	var item *client.Chat
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err = s.clientService.Chat(request.Context(), item)
	if errors.Is(err, client.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetChatStory(writer http.ResponseWriter, request *http.Request) {
	From := request.URL.Query().Get("from")
	To := request.URL.Query().Get("to")
	from, err := strconv.ParseInt(From, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	to, err := strconv.ParseInt(To, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	item, err := s.clientService.ChatStory(request.Context(), from, to)
	if errors.Is(err, client.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
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
