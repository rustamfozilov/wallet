package service

import (
	"github.com/rustamfozilov/wallet/pkg/types"
)

type Service struct {
	accounts *[]types.Account
	payments *[]types.Payment
}

func (a AccountID int64) FindAccountByID (*types.Account, error)  {

}