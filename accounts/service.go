package accounts

import (
	"fmt"
	"sync"
)

type Account struct {
	ID       uint32
	Name     string
	Password string
	Mail     string
}

type Service struct {
	mu sync.RWMutex
	Items []*Account
}

func NewService() *Service {
	return &Service{
		Items: make([]*Account, 0),
	}
}

func (s *Service) Registration(name, password, mail string) (*Account, error){
	s.mu.RLock()
	i := uint32(1)

	for _, item := range s.Items {
		if mail == item.Mail {
			return nil, fmt.Errorf("the mail = %v, is already execute", mail)
		}
	}

	s.Items = append(s.Items, &Account{
		ID:       i,
		Name:     name,
		Password: password,
		Mail:     mail,
	})

	i++

	return s.Items[i], nil
}