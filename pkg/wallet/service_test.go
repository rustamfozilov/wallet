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

	if  !reflect.DeepEqual(*account, want){
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
	if err !=ErrAccountNotFound {
		t.Error(err)
	}


	if  account != nil {
		t.Error(*account)
	}
}