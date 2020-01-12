package types

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	AccountId uint `json:"account_id"`
	AppName string `json:"app_name"`
	AccessUrl string `json:"access_url"`
	RegistryDownloadUrl string `json:"registry_download_url"`
	LocalAccessUrl string `json:"local_access_url"`

	Environments []Environment `json:"environments"`
	Account *Account `json:"account" gorm:"-" sql:"-"`
}

type Environment struct {
	gorm.Model
	AppId uint `json:"app_id"`
	EnvKey string `json:"env_key"`
	EnvValue string `json:"env_value"`
}

func NewEnvVariable(appId uint, k, v string) *Environment {
	return &Environment{
		AppId:    appId,
		EnvKey:   k,
		EnvValue: v,
	}
}

func NewApp(name string, accountId uint) *App {
	return &App{
		AccountId: accountId,
		AppName:   name,
		AccessUrl: fmt.Sprintf("https://%s.hostgoapp.com", name),
	}
}
