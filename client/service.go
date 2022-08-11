package client

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var ErrNotFound = errors.New("account not fount")
var ErrInternal = errors.New("internal error")
var ErrAlreadyExe = errors.New("this account is aleady execute")
var ErrBadRequest = errors.New("you are sending wrong data")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Account struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Chat struct {
	SendlerID  uint64 `json:"sendler"`
	ReceiverID uint64 `json:"receiver"`
	Message    string `json:"message"`
	Time       time.Time
}

var account = &Account{}
var item = &Chat{}

func (s *Service) Registration(ctx context.Context, item *Account) (*Account, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO account(name, password) VALUES($1, $2) RETURNING id, name, password
	`, item.Name, item.Password).Scan(&account.Id, &account.Name, &account.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return account, nil
}

func (s *Service) DeleteAccount(ctx context.Context, item *Account) error {
	row, err := s.pool.Exec(ctx, `
		DELETE FROM account WHERE id = $1
	`, item.Id)
	if err != nil {
		return err
	}
	if row.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	return nil
}

func (s *Service) Chat(ctx context.Context, message *Chat) (*Chat, error) {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO messanger(sendler, receiver, message) VALUES($1, $2, $3) RETURNING sendler, receiver, message
	`, message.SendlerID, message.ReceiverID, message.Message).Scan(&item.SendlerID, &item.ReceiverID, &item.Message)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

func (s *Service) ChatStory(ctx context.Context, sendler, receiver int64) ([]*Chat, error) {
	var data = []*Chat{}
	rows, err := s.pool.Query(ctx, `
		SELECT * FROM messanger WHERE (sendler = $1 AND receiver = $2) OR (sendler = $2 AND receiver = $1)
	`, sendler, receiver)
	if err != nil {
		return nil, err
	}
	log.Println(sendler, receiver)
	defer rows.Close()

	for rows.Next() {
		var item = &Chat{}
		err = rows.Scan(&item.SendlerID, &item.ReceiverID, &item.Message, &item.Time)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return data, nil
}
