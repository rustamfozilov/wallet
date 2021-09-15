package wallet

import (
	"errors"
	"github.com/rustamfozilov/wallet/pkg/types"
)

var ErrAccountNotFound = errors.New("account not found")

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}


func (s Service) FindAccountByID (accountID int64) (*types.Account, error)  {
	for _, account := range s.accounts  {
		if account.ID == accountID{
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}