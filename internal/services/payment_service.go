package services

import (
	"errors"
	"regexp"

	"github.com/LuoZihYuan/Go-Cart/internal/repository"
)

var (
	ErrInvalidPaymentData = errors.New("invalid payment data")
	ErrPaymentDeclined    = errors.New("payment declined")
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	cartRepo    *repository.CartRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository, cartRepo *repository.CartRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		cartRepo:    cartRepo,
	}
}

// ProcessPayment processes a credit card payment
func (s *PaymentService) ProcessPayment(creditCardNumber string, shoppingCartID int) (string, bool, error) {
	// Validate input
	if shoppingCartID < 1 {
		return "", false, ErrInvalidPaymentData
	}

	// Validate credit card number format (13-19 digits)
	matched, _ := regexp.MatchString(`^[0-9]{13,19}$`, creditCardNumber)
	if !matched {
		return "", false, ErrInvalidPaymentData
	}

	// Verify cart exists
	cart, err := s.cartRepo.GetByID(shoppingCartID)
	if err == repository.ErrCartNotFound {
		return "", false, ErrCartNotFound
	}
	if err != nil {
		return "", false, err
	}

	// Calculate total amount (simplified - in real system would use product prices)
	amount := float64(len(cart.Items)) * 10.0 // Mock calculation

	// Process payment
	transactionID, success, err := s.paymentRepo.ProcessPayment(creditCardNumber, shoppingCartID, amount)
	if err != nil {
		return "", false, err
	}

	if !success {
		return transactionID, false, ErrPaymentDeclined
	}

	return transactionID, true, nil
}
