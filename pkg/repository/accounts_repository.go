package repository

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/fn"
	"github.com/saas/hostgolang/pkg/session"
	"github.com/saas/hostgolang/pkg/types"
	"log"
)

var ErrAccountCreateFailed = errors.New("failed to create a new account at this time. Please retry later")
var ErrAccountRetrieveFailed = errors.New("failed to retrieve account account data")
var ErrAuthenticationFailed = errors.New("authentication failed: email and password does not match our record")
var ErrAuthServerFailed = errors.New("failed to finish authentication at this time. Please retry later")

//go:generate mockgen -destination=mocks/accounts_repository_mock.go -package=mocks github.com/saas/hostgolang/pkg/repository AccountRepository
type AccountRepository interface {
	CreateAccount(opt *types.NewAccountOpts) (*types.Account, error)
	AuthenticateAccount(opt *types.AuthenticateAccountOpts) (*types.Account, error)
	GetAccountByEmail(email string) (*types.Account, error)
}

type accountRepository struct {
	db *gorm.DB
	store session.Store
}

func (a *accountRepository) CreateAccount(opt *types.NewAccountOpts) (*types.Account, error) {
	acc, err := a.GetAccountByEmail(opt.Email)
	if err == nil && len(acc.Password) > 0 {
		return nil, errors.New("an account is already linked with that email")
	}
	opt.Password = fn.HashPassword(opt.Password)
	accountToken := fmt.Sprintf("%s:%s", fn.GenerateRandomString(32), fn.GenerateRandomString(36))
	account := types.NewAccount(opt, accountToken)
	tx := a.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	if err := tx.Create(account).Error; err != nil {
		log.Println("failed to create account; ", err)
		tx.Rollback()
		return nil, ErrAccountCreateFailed
	}
	if err := a.store.Put(accountToken, account); err != nil {
		log.Println("failed to create session record; ", err)
		tx.Rollback()
		return nil, ErrAccountCreateFailed
	}
	if err := tx.Commit().Error; err != nil {
		log.Println("failed to commit create account txn", err)
		return nil, ErrAccountCreateFailed
	}
	return account, nil
}

func (a *accountRepository) AuthenticateAccount(opt *types.AuthenticateAccountOpts) (*types.Account, error) {
	account, err := a.GetAccountByEmail(opt.Email)
	if err == gorm.ErrRecordNotFound {
		return nil, ErrAuthenticationFailed
	}
	if err != nil {
		return nil, ErrAuthenticationFailed
	}
	if ok := fn.VerifyHashPassword(account.Password, opt.Password); !ok {
		return nil, ErrAuthenticationFailed
	}
	accountToken := fmt.Sprintf("%s:%s", fn.GenerateRandomString(32), fn.GenerateRandomString(36))
	if err := a.store.Put(accountToken, account); err != nil {
		log.Println("failed to create session record during authentication; ", err)
		return nil, ErrAuthServerFailed
	}
	account.AccountToken = accountToken
	return account, nil
}

func (a *accountRepository) GetAccountByEmail(email string) (*types.Account, error) {
	account := &types.Account{}
	err := a.db.Table("accounts").Where("email = ?", email).First(account).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, err
	case nil:
		return account, nil
	default:
		return nil, ErrAccountRetrieveFailed
	}
}

func NewAccountRepository(db *gorm.DB, store session.Store) AccountRepository {
	return &accountRepository{db:db, store: store}
}
