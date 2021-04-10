package wallet

import "github.com/amirbek-jan/wallet/pkg/types"

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}



func (e Error) Error() string {
	return string(e)
}
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, Error("phone already registered")
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {

}
