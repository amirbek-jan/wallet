package wallet

import "go/types"

type Service struct {
	accounts []types.Account
	payments []types.Payment
}