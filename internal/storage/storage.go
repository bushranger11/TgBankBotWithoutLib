package storage

type Storage struct {
	accounts map[int64]float64
}

func NewStorage() *Storage {
	return &Storage{
		accounts: make(map[int64]float64),
	}
}

func (s *Storage) GetBalance(userID int64) float64 {
	return s.accounts[userID]
}

func (s *Storage) Deposit(userID int64, amount float64) {
	s.accounts[userID] += amount
}

func (s *Storage) Withdraw(userID int64, amount float64) bool {
	if s.accounts[userID] < amount {
		return false
	}
	s.accounts[userID] -= amount
	return true
}
