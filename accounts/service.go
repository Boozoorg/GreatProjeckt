package accounts

import (
	"fmt"
	"sync"
	"time"
)

type Messanger struct{
	From,
	To, 
	Message,
	time string
}

type Account struct {
	ID       uint64
	Name     string
	Password string
}

type Service struct {
	mu sync.RWMutex
	Items []*Account
	message []*Messanger
}

var i uint64 = 1
var j uint64 = 0

func NewService() *Service {
	return &Service{
		Items: make([]*Account, 0),
	}
}

func (s *Service) Registration(name, password string) (*Account, error){
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.Items {
		if name == item.Name {
			return nil, fmt.Errorf("%v is already execute", item.Name)
		}
	}

	s.Items = append(s.Items, &Account{
		ID:       i,
		Name:     name,
		Password: password,
	})

	i++

	return s.Items[i-2], nil
}

func (s *Service)SendMessage(from, to, Message string) (*Messanger, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	a := false
	for _, item := range s.Items {
		if to == item.Name && from == item.Name{
			return nil, fmt.Errorf("you are really send message to your self? o_O")
		}
		if to == item.Name{
			a = true
			
		}
	}

	if !a{
		return nil, fmt.Errorf("there is not any accounts like this %v", to)
	}

	for _, item := range s.Items {
		if from == item.Name {
			s.message = append(s.message, &Messanger{
				From: item.Name,
				To: to,
				Message: Message,
				time: time.Now().Format(time.ANSIC),
			}) 
			j++
			return s.message[j-1], nil
		}
	}

	return nil, fmt.Errorf("there is not any accounts like this %v", from)
}

