package repository

import (
	"errors"
	"github.com/goombaio/namegenerator"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
)

var ErrDuplicateName = errors.New("an app already exists with that name")
var ErrAppNotFound = errors.New("app not found")
var ErrFailedCreateApp = errors.New("failed to create a new app at this time. Please retry later")
var ErrAppUpdateFailed = errors.New("failed to update app. Please retry")
var ErrAppDataRetrieveFailed = errors.New("failed to retrieve app data. Please retry")

//go:generate mockgen -destination=mocks/apps_repository_mock.go -package=mocks github.com/saas/hostgolang/pkg/repository AppsRepository
type AppsRepository interface {
	CreateApp(name string, accountId uint) (*types.App, error)
	GetApp(name string) (*types.App, error)
	GetAccountApps(accountId uint) ([]types.App, error)
	UpdateEnvironmentVars(appName string, vars map[string]string) error
	GetEnvironmentVars(appName string) ([]types.Environment, error)
}

type appsRepository struct {
	db *gorm.DB
	nameGen namegenerator.Generator
	ks8 services.K8sService
}

func (a *appsRepository) GetAccountApps(accountId uint) ([]types.App, error) {
	data := make([]types.App, 0)
	err := a.db.Table("apps").Where("account_id = ?", accountId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *appsRepository) CreateApp(name string, accountId uint) (*types.App, error) {
	if name == "" {
		name = a.nameGen.Generate()
	}
	ap, err := a.GetApp(name)
	if err == nil && len(ap.AppName) > 0 {
		return nil, ErrDuplicateName
	}
	app := types.NewApp(name, accountId)
	if err := a.db.Create(app).Error; err != nil {
		return nil, ErrFailedCreateApp
	}
	return app, nil
}

func (a *appsRepository) GetApp(name string) (*types.App, error) {
	app := &types.App{}
	err := a.db.Table("apps").Where("app_name = ?", name).First(app).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, err
	case nil:
		return app, nil
	default:
		return nil, ErrAppNotFound
	}
}

func (a *appsRepository) UpdateEnvironmentVars(name string, envs map[string]string) error {
	app, err := a.GetApp(name)
	if err != nil {
		return ErrAppNotFound
	}
	appId := app.ID
	m := make([]types.Environment, 0)
	for k, v := range envs {
		if ev, err := a.GetEnvironment(k, appId); err == nil && len(ev.EnvValue) > 0 {
			err = a.DeleteEnvironmentVar(k, appId)
		}
		e := types.NewEnvVariable(appId, k, v)
		if err := a.db.Create(e).Error; err != nil {
			return ErrAppUpdateFailed
		}
		m = append(m, *e)
	}
	return a.ks8.UpdateEnvs(app.AppName, m)
}

func (a *appsRepository) GetEnvironmentVars(name string)([]types.Environment, error) {
	app, err := a.GetApp(name)
	if err != nil {
		return nil, ErrAppNotFound
	}
	data := make([]types.Environment, 0)
	err = a.db.Table("environments").Where("app_id = ?", app.ID).Find(&data).Error
	if err != nil {
		return nil, ErrAppDataRetrieveFailed
	}
	return data, nil
}

func (a *appsRepository) DeleteEnvironmentVar(name string, appId uint) error {
	return a.db.Table("environments").Where("env_key = ? AND app_id = ?", name, appId).Delete(&types.Environment{}).Error
}

func (a *appsRepository) GetEnvironment(name string, appId uint) (*types.Environment, error) {
	env := &types.Environment{}
	err := a.db.Table("environments").Where("env_key = ? AND app_id = ?", name, appId).First(env).Error
	if err != nil {
		return nil, err
	}
	return env, nil
}

func NewAppsRepository(db *gorm.DB, nameGenerator namegenerator.Generator, k8s services.K8sService) AppsRepository {
	return &appsRepository{db:db, nameGen: nameGenerator, ks8: k8s}
}