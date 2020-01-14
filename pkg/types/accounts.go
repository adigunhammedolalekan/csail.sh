package types

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	GithubId       string `json:"github_id"`
	CompanyName    string `json:"company_name"`
	CompanyWebsite string `json:"company_website"`
	AccountToken   string `json:"account_token"`
}

func NewAccount(opt *NewAccountOpts, token string) *Account {
	return &Account{
		Name:           opt.Name,
		Email:          opt.Email,
		Password:       opt.Password,
		CompanyName:    opt.CompanyName,
		CompanyWebsite: opt.CompanyWebsite,
		AccountToken:   token,
	}
}
