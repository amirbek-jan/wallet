package wallet

import (
	"reflect"
	"testing"
	"fmt"
	"github.com/amirbek-jan/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service
}

type testPayment struct {
	ID        string
	AccountID int64
	Amount    types.Money
	Category  types.PaymentCategory
	Status    types.PaymentStatus
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []testPayment
}

var defaultTestPayment = testPayment{
	ID:        uuid.New().String(),
	AccountID: 1,
	Amount:    1_000_00,
	Category:  "auto",
	Status:    types.PaymentStatusInProgress,
}

var defaultTestAccount = testAccount{
	phone:    "+992000000001",
	balance:  10_000_00,
	payments: []testPayment{defaultTestPayment},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't, register account = %v", err)
	}

	// пополняем его счёт
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.Amount, payment.Category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	// регистрируем там пользователя
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счёт
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can,t deposit account, error = %v", err)
	}

	return account, nil
}
func TestService_Reject_success(t *testing.T) {
	// Создаём сервис
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}

}

func TestService_FindPaymentByID_fail(t *testing.T) {

	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID():  must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	repeatPay, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't repeat pay = %v", err)
		return
	}
	if repeatPay.ID == payment.ID {
		t.Errorf("repeat payment ID is not payment ID = %v", repeatPay)
		return
	}
	if repeatPay.Amount != payment.Amount {
		t.Errorf("repeat payment amount is not payment amount = %v", repeatPay)
		return
	}
	if repeatPay.Category != payment.Category {
		t.Errorf("repeat payment category is not payment category = %v", repeatPay)
		return
	}
	if repeatPay.Status != payment.Status {
		t.Errorf("repeat payment status is not payment status = %v", repeatPay)
		return
	}
	if repeatPay.AccountID != payment.AccountID {
		t.Errorf("repeat payment accountID is not payment accountID = %v", repeatPay)
		return
	}
}

func TestService_Repeat_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	if payments != nil {
		return
	}

	repeatPay, err := s.Repeat(uuid.New().String())
	if err == nil {
		t.Errorf("Repeat payment is ok")
		return
	}
	if repeatPay == nil {
		t.Errorf("Repeat payment is ok")
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T){
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "buy")
	if err != nil {
		t.Error(err)
		return
	}

	if favorite.AccountID != payment.AccountID {
		t.Errorf("FavoritePayment(): favorite account ID is not payment account ID")
		return
	}

	if favorite.Categoty != payment.Category {
		t.Errorf("FavoritePayment(): favorite category is not payment category")
		return
	}

	if favorite.Amount != payment.Amount {
		t.Errorf("FavoritePayment(): favorite amount is not payment amount")
		return
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	if payments != nil {
		return
	}

	favorite, err := s.FavoritePayment("123456789", "buy")
	if err == nil {
		t.Errorf("dbsthtsnts")
		return
	}
	if favorite == nil {
		t.Errorf("Favorite payment is ok")
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "buy")
	if err != nil {
		t.Error(err)
		return
	}

	payFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if payFavorite.AccountID != favorite.AccountID {
		t.Errorf("FavoritePayment(): favorite account ID is not payment account ID")
		return
	}
	if payFavorite.Amount != favorite.Amount {
		t.Errorf("FavoritePayment(): favorite amount is not payment amount")
		return
	}
	if payFavorite.Category != favorite.Categoty {
		t.Errorf("FavoritePayment(): favorite category is not payment category")
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "buy")
	if err != nil {
		t.Error(err)
		return
	}

	payFavorite, err := s.PayFromFavorite("asdasdasd")
	if err == nil {
		t.Error(err)
		return
	}

	if payFavorite != nil {
		t.Errorf("invalid favorite")
		return
	}

}