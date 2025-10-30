package repository

import (
	"fmt"
	"sync"
	"time"
)

type PaymentTransaction struct {
	TransactionID    string
	CreditCardNumber string
	ShoppingCartID   int
	Amount           float64
	Timestamp        time.Time
	Success          bool
}

type PaymentRepository struct {
	transactions map[string]*PaymentTransaction
	mu           sync.RWMutex
	nextTxID     int
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		transactions: make(map[string]*PaymentTransaction),
		nextTxID:     1,
	}
}

// ProcessPayment processes a payment and stores the transaction
func (r *PaymentRepository) ProcessPayment(creditCardNumber string, shoppingCartID int, amount float64) (string, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simple validation: card numbers ending in 0 are declined
	lastDigit := creditCardNumber[len(creditCardNumber)-1]
	success := lastDigit != '0'

	transactionID := fmt.Sprintf("txn_%d_%d", time.Now().Unix(), r.nextTxID)
	r.nextTxID++

	transaction := &PaymentTransaction{
		TransactionID:    transactionID,
		CreditCardNumber: creditCardNumber,
		ShoppingCartID:   shoppingCartID,
		Amount:           amount,
		Timestamp:        time.Now(),
		Success:          success,
	}

	r.transactions[transactionID] = transaction

	return transactionID, success, nil
}

// GetTransaction retrieves a transaction by ID
func (r *PaymentRepository) GetTransaction(transactionID string) (*PaymentTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, exists := r.transactions[transactionID]
	if !exists {
		return nil, ErrProductNotFound // Reusing error, could create ErrTransactionNotFound
	}

	txCopy := *tx
	return &txCopy, nil
}
