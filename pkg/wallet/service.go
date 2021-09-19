package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rustamfozilov/wallet/pkg/types"
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

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error)  {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	var favorite = types.Favorite{
		ID:        payment.ID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, &favorite)
	return &favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error)  {
	var targetFavorite *types.Favorite
	for _, favorite := range s.favorites {
		if favorite.ID== favoriteID{
		targetFavorite = favorite
		break
		}
		return nil, ErrFavoriteNotFound
	}
	payment := &types.Payment{
		ID:        uuid.New().String(),
		AccountID: targetFavorite.AccountID,
		Amount:    targetFavorite.Amount,
		Category:  targetFavorite.Category,
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