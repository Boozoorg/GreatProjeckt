package accounts

import (
	"fmt"
	"sync"
)

var n []uint64

type Account struct {
	ID       uint64
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
	i := uint64(1)

	for _, item := range s.Items {
		if mail == item.Mail {
			return nil, fmt.Errorf("the mail = %v, is already execute", mail)
		}
	}

	for _, item := range n {
		if n != nil{
			break
		}
		if item != 0 {
			s.Items[item] = &Account{
				ID:       item,
				Name:     name,
				Password: password,
				Mail:     mail,
			}
			n[item] = 0
				
			return s.Items[i], nil
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

func (s *Service) DelateAccByID(id uint64) (*Account, error){
	for _, item := range s.Items {
		
		if id == item.ID {
			s.Items[id] = nil
			n = append(n, id)
			return s.Items[id], nil
		}
	}

	return nil, fmt.Errorf("this account is not exist")
}

