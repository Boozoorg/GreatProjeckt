package client

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAlreadyExe      = errors.New("this account is aleady execute")
	ErrBadRequest      = errors.New("you are sending wrong data")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoRow           = errors.New("no row found to delete")
	ErrNotFound        = errors.New("account not fount")
	ErrInternal        = errors.New("internal error")
	ErrNoSuchUser      = errors.New("no such user")
)

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

type Token struct {
	Id    uint64 `json:"id"`
	Token string `json:"token"`
}

type Chat struct {
	SendlerID  uint64    `json:"sendler"`
	ReceiverID uint64    `json:"receiver"`
	Message    string    `json:"message"`
	Time       time.Time `json:"time"`
}

var account = &Account{}
var item = &Chat{}

func (s *Service) Registration(ctx context.Context, item *Account) (*Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	err = s.pool.QueryRow(ctx, `
		INSERT INTO account(name, password) VALUES($1, $2) RETURNING id, name, password
	`, item.Name, string(hash)).Scan(&account.Id, &account.Name, &account.Password)
	if err != nil {
		log.Print(err)
		return nil, ErrAlreadyExe
	}

	return &Account{Id: account.Id, Name: account.Name, Password: item.Password}, nil
}

func (s *Service) TokenToClient(ctx context.Context, item *Account) (*Token, error) {
	var hash string
	err := s.pool.QueryRow(ctx, `
		SELECT id, password FROM account WHERE name = $1
	`, &item.Name).Scan(&item.Id, &hash)
	if err == pgx.ErrNoRows {
		log.Print(err)
		return nil, ErrNoSuchUser
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(item.Password))
	if err != nil {
		log.Print(err)
		return nil, ErrNoSuchUser
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `
		INSERT INTO account_token(token, account_id) VALUES($1, $2)
	`, token, &item.Id)
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return &Token{Token: token, Id: item.Id}, nil
}

func (s *Service) IDFunc(ctx context.Context, token string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
		SELECT account_id FROM account_token WHERE token = $1
	`, token).Scan(&id)

	if err == pgx.ErrNoRows {
		return 0, ErrNoRow
	}
	if err != nil {
		log.Print(err)
		return 0, ErrInternal
	}

	return id, nil
}

func (s *Service) DeleteAccount(ctx context.Context, id int64) error {
	row, err := s.pool.Exec(ctx, `
		DELETE FROM account WHERE id = $1
	`, id)
	if err != nil {
		log.Print(err)
		return err
	}
	if row.RowsAffected() != 1 {
		log.Print(err)
		return ErrNoRow
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
		log.Print(err)
		return nil, err
	}
	log.Println(sendler, receiver)
	defer rows.Close()

	for rows.Next() {
		var item = &Chat{}
		err = rows.Scan(&item.SendlerID, &item.ReceiverID, &item.Message, &item.Time)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		data = append(data, item)
	}

	if rows.Err() != nil {
		log.Print(err)
		return nil, err
	}

	return data, nil
}
