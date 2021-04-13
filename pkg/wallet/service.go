package wallet

import (
	"errors"
	
	"github.com/amirbek-jan/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrAccountNotFound = errors.New("account not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites	  []*types.Favorite	
}

type Error string

func (e Error) Error() string {
	return string(e)
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

// Find Payment by ID
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

// Find Account by ID
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrPaymentNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)
	account.Balance -= payment.Amount
	newPaymentID := uuid.New().String()
	newPayment := &types.Payment{
		ID: newPaymentID,
		AccountID: payment.AccountID,
		Amount: payment.Amount,
		Category: payment.Category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, newPayment)
	return newPayment,nil
	
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	newFavoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:       newFavoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Categoty:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	targetFavorite := &types.Favorite{}
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			targetFavorite = favorite
		}
	}

	if targetFavorite == nil {
		return nil, ErrFavoriteNotFound
	}

	account, err := s.FindAccountByID(targetFavorite.AccountID)
	if err != nil {
		return nil, ErrFavoriteNotFound
	}
	account.Balance -= targetFavorite.Amount
	newPaymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        newPaymentID,
		AccountID: targetFavorite.AccountID,
		Amount:    targetFavorite.Amount,
		Category:  targetFavorite.Categoty,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}
