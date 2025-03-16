package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(connString string) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetBalance(userID int64) (float64, error) {
	var balance float64
	err := s.db.QueryRow(context.Background(), "SELECT balance FROM balances WHERE user_id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении баланса: %w", err)
	}
	return balance, nil
}

func (s *Storage) Deposit(userID int64, amount float64) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO balances (user_id, balance)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET balance = balances.balance + EXCLUDED.balance
	`, userID, amount)
	if err != nil {
		return fmt.Errorf("ошибка при пополнении баланса: %w", err)
	}
	return nil
}

func (s *Storage) Withdraw(userID int64, amount float64) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE balances
		SET balance = balance - $1
		WHERE user_id = $2 AND balance >= $1
	`, amount, userID)
	if err != nil {
		return fmt.Errorf("ошибка при снятии средств: %w", err)
	}
	return nil
}
