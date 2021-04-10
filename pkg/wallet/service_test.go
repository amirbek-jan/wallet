package wallet

import (
	"testing"

	"github.com/amirbek-jan/wallet/pkg/types"
)

// TestService func
func TestService_Reject_PaymentStatusFail(t *testing.T) {
	var service Service
	payments := []*types.Payment{
		{
			"1",
			1,
			1000,
			"mobile",
			types.PaymentStatusInProgress,
		},
		{
			"2",
			1,
			2000,
			"device",
			types.PaymentStatusInProgress,
		},
		{
			"3",
			1,
			3000,
			"transport",
			types.PaymentStatusFail,
		},
	}
	service.payments = payments
	payment, err := service.FindPaymentByID("4")
	if err == nil && payment != nil {
		t.Error(err)
	}
}

func TestService_Reject_PaymentStatusFound(t *testing.T) {
	var service Service
	payments := []*types.Payment{
		{
			"1",
			1,
			1000,
			"mobile",
			types.PaymentStatusInProgress,
		},
		{
			"2",
			1,
			2000,
			"device",
			types.PaymentStatusInProgress,
		},
		{
			"3",
			1,
			3000,
			"transport",
			types.PaymentStatusFail,
		},
	}
	service.payments = payments
	payment, err := service.FindPaymentByID("1")
	if payment == nil && err != nil {
		t.Error(err)
	}
}
