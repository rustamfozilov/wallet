package wallet

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/rustamfozilov/wallet/pkg/types"
	"log"
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

func TestService_FindPaymentByID_successs(t *testing.T) {
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

	if !reflect.DeepEqual(*payment, want) {
		t.Errorf("invalid payment : %v, want : %v", payment, want)
	}

}

func TestService_FindPaymentByID_faill(t *testing.T) {
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

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("cant register acoount , error := %v ", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("cant deposite account, error := %v ", err)
	}
	return account, nil
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultAccount = testAccount{
	phone:   "+992985410248",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register acoount , error := %v ", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposite account, error := %v ", err)
	}
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant make payment, error := %v ", err)
		}
	}
	return account, payments, nil
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Fatal(err)
	}
	payment := payments[0]

	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID() : error = %v ", err)
		return
	}
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Fatal(err)
	}


	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID() : must return error returned nil ")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("ne pechataet")
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Error(err)
	}
	payment := payments[0]
	repeatedPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(payment,repeatedPayment) {
			t.Errorf("not change paimentID payment : %v, repeated paymentID %v",payment.ID, repeatedPayment.ID)}

log.Println(*payment, *repeatedPayment)
}

func TestService_Repeat_fail(t *testing.T) {
	s := newTestService()

	_, _, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Error(err)
	}

	_, err = s.Repeat(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID() : must return error returned nil ")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("ne pechataet")
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Fatal(err)
	}
	payment := payments[0]
	var  testFavitePayment = types.Favorite{
		ID:        "e6b3cf42-3deb-429b-978c-e93165707148",
		AccountID: payment.AccountID,
		Name:      "Hahah",
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	got, err := s.FavoritePayment(payment.ID, "Hahah")
	if err != nil {
		t.Fatal(err)
	}
	if testFavitePayment.ID == got.ID {
		t.Errorf("invalid made favorite payment want: %v, got: %v", testFavitePayment,got)
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultAccount)
	if err != nil {
		t.Fatal(err)
	}
	payment := payments[0]
	var  testPayment = types.Payment{
		ID:        payment.ID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	fav, err := s.FavoritePayment(payment.ID, "Hahah")
	if err != nil {
		t.Fatal(err)
	}
	payFromFavorite, err := s.PayFromFavorite(fav.ID)
	if err != nil {
		t.Fatal(err)
	}
if testPayment.ID == payFromFavorite.ID {
	t.Fatal("id not changed/ want:",testPayment.ID,"got :", payFromFavorite.ID)
}
}

func Test_readFromFileInSliceByte(t *testing.T) {
	var b bytes.Buffer
	s := "hello world"
	b.WriteString(s)
	a, err := readAll(&b)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(a) {
		t.Fatal(s, string(a))
	}
	b.Reset()

	s = ""
	b.WriteString(s)
	a, err = readAll(&b)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(a) {
		t.Fatal(s, string(a))
	}
}
