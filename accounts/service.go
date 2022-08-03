package accounts

import (
	"fmt"
	"sync"
)

type key struct {
	From,
	To string
}

type Messanger struct {
	Message,
	Time string
}

var Message = make(map[*key][]*Messanger)

type Account struct {
	ID       uint64
	Name     string
	Password string
}

type Service struct {
	mu      sync.RWMutex
	Items   []*Account
}

var i uint64 = 1

func NewService() *Service {
	return &Service{
		Items: make([]*Account, 0),
	}
}

func (s *Service) Registration(name, password string) (*Account, error) {
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

func (s *Service) SendMessage(from, to, message, time string) (*Messanger, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	a := false
	for _, item := range s.Items {
		if to == item.Name && from == item.Name {
			return nil, fmt.Errorf("you are really send message to your self? o_O")
		}
		if to == item.Name {
			a = true
		}
	}

	if !a {
		return nil, fmt.Errorf("there is not any accounts like this %v", to)
	}

	for _, item := range s.Items {
		if from == item.Name {
			Message[&key{From: from, To: to}] = append(Message[&key{From: from, To: to}], &Messanger{
				Message: message,
				Time:    time,
			})

			return Message[&key{From: from, To: to}][len(Message[&key{From: from, To: to}])-1], nil
		}
	}

	return nil, fmt.Errorf("there is not any accounts like this %v", from)
}

func (s *Service) MessageStory(Person1, Person2 string){
	s.mu.RLock()
	defer s.mu.RUnlock()

	// a := false
	// for _, item := range s.Items {
	// 	if Person1 == item.Name && Person2 == item.Name {
	// 		return nil, fmt.Errorf("you are really want to read message that you sand to your self? O_o")
	// 	}
	// 	if Person1 == item.Name {
	// 		a = true
	// 	}
	// }

	// if !a {
	// 	return nil, fmt.Errorf("there is not any accounts like this %v", Person1)
	// }

	// for _, item := range s.Items {
	// 	if Person2 == item.Name {
	// 		if Message[&key{From: Person1, To: Person2}] == nil || Message[&key{From: Person2, To: Person1}] == nil {
	// 			return nil, "you did not chated yet"
	// 		}

	// 		a := 0
	// 		for _, v1 := range Message[&key{From: Person1, To: Person2}] {
	// 			for key, v2 := range Message[&key{From: Person2, To: Person1}] {
	// 				if v1.Time < v2.Time{
	// 					s.message = append(s.message, v1)
	// 					a = key
	// 					break
	// 				}
					
	// 				s.message = append(s.message, v2)
	// 			}
	// 		}
	// 	}
	// }

	// j := 0
	// for _, item = range s.Items {
	// 	fmt.Println(Message[&key{From: Person1, To: Person1}][j])
	// 	j++
	// }

	// return nil, fmt.Errorf("there is not any accounts like this %v", Person2)
}