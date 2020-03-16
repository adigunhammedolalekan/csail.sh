package repository

import (
	"errors"
	"github.com/goombaio/namegenerator"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"log"
)

var ErrDuplicateName = errors.New("an app already exists with that name")
var ErrAppNotFound = errors.New("app not found")
var ErrFailedCreateApp = errors.New("failed to create a new app at this time. Please retry later")
var ErrAppUpdateFailed = errors.New("failed to update app. Please retry")
var ErrAppDataRetrieveFailed = errors.New("failed to retrieve app data. Please retry")
var ErrDomainNotFound = errors.New("domain not found")
var ErrDomainExists = errors.New("domain already exists in the system")

//go:generate mockgen -destination=mocks/apps_repository_mock.go -package=mocks github.com/saas/hostgolang/pkg/repository AppsRepository
type AppsRepository interface {
	CreateApp(name, plan string, accountId uint) (*types.App, error)
	GetApp(name string) (*types.App, error)
	GetAppById(appId uint) (*types.App, error)
	GetAccountApps(accountId uint) ([]types.App, error)
	UpdateEnvironmentVars(appName string, vars map[string]string) error
	GetEnvironmentVars(appName string) ([]types.Environment, error)
	DeleteEnvironmentVars(appName string, keys []string) error
	ScaleApp(appName string, replicas int) error
	ListRunningInstances(appName string) ([]types.Instance, error)
	GetPlan(appId uint) (*types.Plan, error)
	UpdatePlan(appId uint, planAlias string) error
	GetDomain(address string) (*types.Domain, error)
	GetDomainByAppId(appId uint) (*types.Domain, error)
	CreateDomain(appId uint, address string) (*types.Domain, error)
	RemoveDomain(appId uint, address string) error
}

type appsRepository struct {
	db      *gorm.DB
	nameGen namegenerator.Generator
	ks8     services.K8sService
	gitService services.GitService
}

func (a *appsRepository) GetAccountApps(accountId uint) ([]types.App, error) {
	data := make([]types.App, 0)
	err := a.db.Table("apps").Where("account_id = ?", accountId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *appsRepository) CreateApp(name, plan string, accountId uint) (*types.App, error) {
	if name == "" {
		name = a.nameGen.Generate()
	}
	ap, err := a.GetApp(name)
	if err == nil && len(ap.AppName) > 0 {
		return nil, ErrDuplicateName
	}
	tx := a.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	app := types.NewApp(name, accountId)
	if err := tx.Create(app).Error; err != nil {
		return nil, ErrFailedCreateApp
	}
	appPlan := types.NewPlan(app.ID, plan)
	if err := tx.Create(appPlan).Error; err != nil {
		tx.Rollback()
		return nil, ErrFailedCreateApp
	}
	if err := a.gitService.CreateRepository(name); err != nil {
		log.Println("failed to create repository ", err)
		tx.Rollback()
		return nil, ErrFailedCreateApp
	}
	if err := tx.Commit().Error; err != nil {
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

func (a *appsRepository) GetAppById(appId uint) (*types.App, error) {
	app := &types.App{}
	err := a.db.Table("apps").Where("id = ?", appId).First(app).Error
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
	for k, v := range envs {
		if ev, err := a.GetEnvironment(k, appId); err == nil && len(ev.EnvValue) > 0 {
			err = a.DeleteEnvironmentVar(k, appId)
		}
		e := types.NewEnvVariable(appId, 0, k, v)
		if err := a.db.Create(e).Error; err != nil {
			return ErrAppUpdateFailed
		}
	}
	m, err := a.GetEnvironmentVars(app.AppName)
	if err != nil {
		return err
	}
	return a.ks8.UpdateEnvs(app.AppName, m)
}

func (a *appsRepository) GetEnvironmentVars(name string) ([]types.Environment, error) {
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

func (a *appsRepository) DeleteEnvironmentVars(appName string, keys []string) error {
	app, err := a.GetApp(appName)
	if err != nil {
		return ErrAppNotFound
	}
	for _, key := range keys {
		if err := a.DeleteEnvironmentVar(key, app.ID); err != nil {
			return ErrAppUpdateFailed
		}
	}
	newEnvs, err := a.GetEnvironmentVars(appName)
	if err != nil {
		return err
	}
	return a.ks8.UpdateEnvs(appName, newEnvs)
}

func (a *appsRepository) ScaleApp(appName string, replicas int) error {
	app, err := a.GetApp(appName)
	if err != nil {
		return ErrAppNotFound
	}
	err = a.db.Table("deployment_settings").Where("app_id = ?", app.ID).
		UpdateColumn("replicas", replicas).Error
	if err != nil {
		return ErrAppUpdateFailed
	}
	return a.ks8.ScaleApp(appName, replicas)
}

func (a *appsRepository) ListRunningInstances(appName string) ([]types.Instance, error) {
	return a.ks8.ListRunningPods(appName)
}

func (a *appsRepository) GetPlan(appId uint) (*types.Plan, error) {
	p := &types.Plan{}
	err := a.db.Table("plans").Where("app_id = ?", appId).First(p).Error
	if err != nil {
		return &types.DefaultPlan, nil
	}
	return p, nil
}

func (a *appsRepository) UpdatePlan(appId uint, planAlias string) error {
	app, err := a.GetAppById(appId)
	if err != nil {
		return err
	}
	_ = app
	return errors.New("not yet implemented")
}

func (a *appsRepository) GetDomain(address string) (*types.Domain, error) {
	d := &types.Domain{}
	err := a.db.Table("domains").Where("address = ?", address).First(d).Error
	if err != nil {
		return nil, ErrDomainNotFound
	}
	return d, nil
}

func (a *appsRepository) GetDomainByAppId(appId uint) (*types.Domain, error) {
	d := &types.Domain{}
	err := a.db.Table("domains").Where("app_id = ?", appId).First(d).Error
	if err != nil {
		return nil, ErrDomainNotFound
	}
	return d, nil
}

func (a *appsRepository) CreateDomain(appId uint, address string) (*types.Domain, error) {
	app, err := a.GetAppById(appId)
	if err != nil {
		return nil, err
	}
	dom, err := a.GetDomain(address)
	if err == nil && dom.Address != "" {
		return nil, ErrDomainExists
	}
	err = a.ks8.AddDomain(app.AppName, address)
	if err != nil {
		return nil, err
	}
	dom = types.NewDomain(appId, address)
	if err := a.db.Create(dom).Error; err != nil {
		return nil, err
	}
	return dom, nil
}

func (a *appsRepository) RemoveDomain(appId uint, address string) error {
	app, err := a.GetAppById(appId)
	if err != nil {
		return err
	}
	err = a.db.Table("domains").Where("app_id = ? AND address = ?", appId, address).Delete(&types.Domain{}).Error
	if err != nil {
		return err
	}
	if err := a.ks8.RemoveDomain(app.AppName, address); err != nil {
		return err
	}
	return nil
}

func NewAppsRepository(db *gorm.DB, nameGenerator namegenerator.Generator,
	k8s services.K8sService, gitService services.GitService) AppsRepository {
	return &appsRepository{db: db, nameGen: nameGenerator, ks8: k8s, gitService: gitService}
}
