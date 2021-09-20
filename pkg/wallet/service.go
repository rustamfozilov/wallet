package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rustamfozilov/wallet/pkg/types"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

//var ErrAccountNotFound = errors.New("account not found")
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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
		return ErrAccountNotFound
	}
	account.Balance += amount
	return nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

//func (s *Service) Reject(paymentID string) error {
//	for _, payment := range s.payments {
//		if payment.ID == paymentID {
//			payment.Status = types.PaymentStatusFail
//			_, err := s.FindAccountByID(payment.AccountID)
//			if err != nil {
//				return err
//			}
//			//acc.Balance = acc.Balance - payment.Amount
//			return nil
//		}
//	}
//	return ErrPaymentNotFound
//}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	//log.Print("[")
	//for _, payment := range s.payments {
	//	log.Print(*payment)
	//}
	//log.Println("]")
	//log.Println("paymentID:", paymentID)
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	payment := &types.Payment{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	// to do acc
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}
	account.Balance = account.Balance - payment.Amount
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s Service) Reject(paymentID string) error {
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

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}
	var repeatedPayment = types.Payment{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    payment.Status,
	}
	//log.Println("reapetedPayment",repeatedPayment)
	s.payments = append(s.payments, &repeatedPayment)
	account.Balance = account.Balance - payment.Amount
	return &repeatedPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	//log.Print("[")
	//for _, payment := range s.payments {
	//	log.Print(*payment)
	//}
	//log.Println("]")
	//log.Println("paymentID:", paymentID)
	var favorite = types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, &favorite)
	return &favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {

	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

func (s *
Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, account := range s.accounts {
		line := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10) + "|"
		_, err = file.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	accounts, err2 := readAll(file)
	if err2 != nil {
		return err2
	}
	lines := strings.Split(string(accounts), "|")
	for _, line := range lines {
		fields := strings.Split(line, ";")
		idString := fields[0]
		id, err := strconv.ParseInt(idString, 10, 64)

		if err != nil {
			return err
		}

		balanceString := fields[1]
		balance, err := strconv.ParseInt(balanceString, 10, 64)
		if err != nil {
			return err
		}
		var acc = types.Account{
			ID:id,
			Phone:   types.Phone(fields[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, &acc)
	}
	return nil
}


func readAll(reader io.Reader) ([]byte, error) {
	accounts := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		read, err := reader.Read(buf)
		if err == io.EOF {
			accounts = append(accounts, buf[:read]...)
			break
		}
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, buf[:read]...)
	}
	return accounts, nil
}
