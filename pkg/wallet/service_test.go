package wallet

import (
	"github.com/rustamfozilov/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestService_FindAccountByID(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 0,
			},
		},
		payments: nil,
	}

	account, err := service.FindAccountByID(1)
	if err != nil {
		t.Error(err)
	}
	want := types.Account{
		ID:      1,
		Phone:   "123",
		Balance: 0,
	}

	if !reflect.DeepEqual(*account, want) {
		t.Error("want:", want, "got:", *account)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 0,
			},
		},
		payments: nil,
	}

	account, err := service.FindAccountByID(3)
	if err != ErrAccountNotFound {
		t.Error(err)
	}

	if account != nil {
		t.Error(*account)
	}
}

func TestServiceS_reject_success(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 10,
			},
		},
		payments: []*types.Payment{
			{
				ID:        "123",
				AccountID: 2,
				Amount:    20,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
			{
				ID:        "124",
				AccountID: 1,
				Amount:    40,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
		},
	}

	err := service.Reject("123")
	if err != nil {
		t.Error(err)
	}
	payment, err := service.FindPaymentByID("123")
	if err != nil {
		t.Error(err)
	}

	wantPayment:= types.Payment{
		ID:        "123",
		AccountID: 2,
		Amount:    20,
		Category:  "",
		Status:    types.PaymentStatusFail,
	}
	if !reflect.DeepEqual(*payment,wantPayment) {
		t.Errorf("invalid changed payment: %v, want : %v ", payment, wantPayment)
	}

	account, err := service.FindAccountByID(2)
	if err != nil {
		t.Error(err)
	}
	wantAccount := types.Account{
		ID:      2,
		Phone:   "321",
		Balance: 30,
	}
	if !reflect.DeepEqual(*account,wantAccount) {
		t.Errorf("invalid changed account: %v, want : %v ", *account, wantAccount)
	}
}


func TestServiceS_reject_fail(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 10,
			},
		},
		payments: []*types.Payment{
			{
				ID:        "123",
				AccountID: 2,
				Amount:    20,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
			{
				ID:        "124",
				AccountID: 1,
				Amount:    40,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
		},
	}

	err := service.Reject("126")
	if err != ErrPaymentNotFound {
		t.Error(err)
	}




}

func TestService_FindPaymentByID_success(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 10,
			},
		},
		payments: []*types.Payment{
			{
				ID:        "123",
				AccountID: 2,
				Amount:    20,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
			{
				ID:        "124",
				AccountID: 1,
				Amount:    40,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
		},
	}
	payment, err := service.FindPaymentByID("123")
	if err != nil {
		t.Error(err)
	}
	want := types.Payment{
			ID:        "123",
			AccountID: 2,
			Amount:    20,
			Category:  "",
			Status:    types.PaymentStatusInProgress,
		}


	if !reflect.DeepEqual(*payment,want){
		t.Errorf("invalid payment : %v, want : %v", payment,want)
	}

}





func TestService_FindPaymentByID_fail(t *testing.T) {
	service := Service{
		nextAccountID: 0,
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "123",
				Balance: 0,
			},
			{
				ID:      2,
				Phone:   "321",
				Balance: 10,
			},
		},
		payments: []*types.Payment{
			{
				ID:        "123",
				AccountID: 2,
				Amount:    20,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
			{
				ID:        "124",
				AccountID: 1,
				Amount:    40,
				Category:  "",
				Status:    types.PaymentStatusInProgress,
			},
		},
	}
	payment, err := service.FindPaymentByID("54664")
	if err != ErrPaymentNotFound {
		t.Error(err)
	}
	if payment != nil {
		t.Error(*payment)
	}
}