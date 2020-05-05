package repository

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/config"
	proxy "github.com/saas/hostgolang/pkg/proxyclient"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"log"
)

var ErrDeploymentFailed = errors.New("deployment failed. Please contact support")
var ErrNoChangeToDeploy = errors.New("no changes to deploy")
var ErrNotFound = errors.New("resources not found")
var ErrReleaseNotFound = errors.New("release not found")

const tempClonePath = "/tmp/git"
//go:generate mockgen -destination=../mocks/deployment_repository_mock.go -package=mocks github.com/saas/hostgolang/pkg/repository DeploymentRepository
type DeploymentRepository interface {
	CreateDeployment(app *types.App, dockerUrl string) (*types.DeploymentResult, error)
	UpdateEnvironmentVars(app *types.App, envs []types.Environment) error
	GetApplicationLogs(appName string) (string, error)
	GetRelease(appId uint) (*types.Release, error)
	GetReleaseByVersion(appId uint, version string) (*types.Release, error)
	CreateRelease(app *types.App, ref string) error
	CreateOrUpdateDeploymentSettings(appId, replicas uint) error
	GetDeploymentSettings(appId uint) (*types.DeploymentSettings, error)
	RollbackDeployment(appId uint, version string) (*types.DeploymentResult, error)
	GetReleases(appId uint) ([]types.Release, error)
	HasRegistryAuthorization(req *types.AuthorizationRequest) ([]string, error)
}

type defaultDeploymentRepository struct {
	db         *gorm.DB
	docker     services.DockerService
	k8s        services.K8sService
	proxy      proxy.Client
	appRepo    AppsRepository
	storage    services.StorageClient
	cfg        *config.Config
}

func (d *defaultDeploymentRepository) CreateDeployment(app *types.App, dockerUrl string) (*types.DeploymentResult, error) {
	rls, err := d.GetRelease(app.ID)
	if err == nil && rls.DockerUrl == dockerUrl {
		return nil, ErrNoChangeToDeploy
	}
	if err := d.CreateRelease(app, dockerUrl); err != nil {
		return nil, err
	}

	appName := app.AppName
	var replicas uint = 1
	settings, err := d.GetDeploymentSettings(app.ID)
	if err == nil {
		replicas = settings.Replicas
	}

	envs, _ := d.appRepo.GetEnvironmentVars(appName)
	m := make(map[string]string, 0)
	if envs != nil && len(envs) > 0 {
		for _, e := range envs {
			m[e.EnvKey] = e.EnvValue
		}
	}
	opt := &types.CreateDeploymentOpts{
		Envs:     m,
		Name:     appName,
		Replicas: int32(replicas),
		Tag:      dockerUrl,
		IsLocal:  true,
		Memory:   settings.Plan.Info.Memory,
		Cpu:      settings.Plan.Info.Cpu,
	}
	result, err := d.k8s.DeployService(opt)
	if err != nil {
		log.Println("failed to deploy service: ", err)
		return nil, ErrDeploymentFailed
	}
	if err := d.proxy.Set(appName, result.Address); err != nil {
		log.Println("failed to contact proxy server: ", err)
		return nil, ErrDeploymentFailed
	}

	if err := d.CreateOrUpdateDeploymentSettings(app.ID, replicas); err != nil {
		log.Println("failed to create deployment settings: ", err)
	}
	result.Version = fmt.Sprintf("v%d", rls.VersionNumber + 1)
	result.Address = fmt.Sprintf("https://%s.%s", app.AppName, d.cfg.ServerUrl)
	return result, nil
}

func (d *defaultDeploymentRepository) UpdateEnvironmentVars(app *types.App, envs []types.Environment) error {
	return d.k8s.UpdateEnvs(app.AppName, envs)
}

func (d *defaultDeploymentRepository) GetApplicationLogs(appName string) (string, error) {
	return d.k8s.GetLogs(appName)
}

func (d *defaultDeploymentRepository) GetRelease(appId uint) (*types.Release, error) {
	r := &types.Release{}
	err := d.db.Table("releases").Where("app_id = ?", appId).Last(r).Error
	if err != nil {
		return r, ErrReleaseNotFound
	}
	return r, nil
}

func (d *defaultDeploymentRepository) updateRelease(r *types.Release) error {
	newVersionNumber := r.VersionNumber + 1
	r.VersionNumber = newVersionNumber
	return d.db.Table("releases").Where("id = ?", r.ID).Update(r).Error
}

func (d *defaultDeploymentRepository) CreateRelease(app *types.App, ref string) error {
	release, err := d.GetRelease(app.ID)
	if err != nil {
		// a fresh app. create a new release
		rls := types.NewRelease(app.ID, ref, 1)
		return d.db.Create(rls).Error
	}
	// use last stable release if wasn't a deploy
	if ref == "" {
		ref = release.DockerUrl
	}
	envs, _ := d.appRepo.GetEnvironmentVars(app.AppName)
	if envs == nil {
		envs = make([]types.Environment, 0)
	}
	newReleaseNumber := release.VersionNumber + 1
	version := fmt.Sprintf("v%d", newReleaseNumber)
	key := fmt.Sprintf("%s:%s", app.AppName, version)
	rlsConfig := types.NewReleaseConfig(version, envs)
	if err := d.storage.PutReleaseConfig(key, rlsConfig); err != nil {
		return err
	}
	rls := types.NewRelease(app.ID, ref, newReleaseNumber)
	return d.db.Create(rls).Error
}

func (d *defaultDeploymentRepository) GetReleaseByVersion(appId uint, version string) (*types.Release, error) {
	r := &types.Release{}
	err := d.db.Table("releases").Where("app_id = ? AND version_number = ?", appId, version).First(r).Error
	if err != nil {
		return r, ErrReleaseNotFound
	}
	return r, nil
}

func (d *defaultDeploymentRepository) GetReleases(appId uint) ([]types.Release, error) {
	data := make([]types.Release, 0)
	err := d.db.Table("releases").Where("app_id = ?", appId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *defaultDeploymentRepository) CreateOrUpdateDeploymentSettings(appId, replicas uint) error {
	if _, err := d.GetDeploymentSettings(appId); err != nil {
		settings := types.NewDeploymentSettings(appId, replicas)
		return d.db.Create(settings).Error
	}
	return d.db.Table("deployment_settings").Where("app_id = ?", appId).UpdateColumn("replicas", replicas).Error
}

func (d *defaultDeploymentRepository) GetDeploymentSettings(appId uint) (*types.DeploymentSettings, error) {
	settings := &types.DeploymentSettings{
		Plan: &types.DefaultPlan,
	}
	err := d.db.Table("deployment_settings").Where("app_id = ?", appId).First(settings).Error
	if err != nil {
		return settings, ErrNotFound
	}
	p, _ := d.appRepo.GetPlan(appId)
	settings.Plan = p
	return settings, nil
}

func (d *defaultDeploymentRepository) RollbackDeployment(appId uint, version string) (*types.DeploymentResult, error) {
	currentRelease, err := d.GetRelease(appId)
	if err != nil {
		return nil, err
	}
	if fmt.Sprintf("v%d", currentRelease.VersionNumber) == version {
		return nil, errors.New("cannot rollback to current release")
	}
	versionNumber := version[1:]
	rls, err := d.GetReleaseByVersion(appId, versionNumber)
	if err != nil {
		return nil, err
	}
	app, err := d.appRepo.GetAppById(appId)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s:%s", app.AppName, version)
	cfg, err := d.storage.GetReleaseConfig(key)
	if err != nil {
		return nil, err
	}
	appName := app.AppName
	var replicas uint = 1
	settings, err := d.GetDeploymentSettings(app.ID)
	if err == nil {
		replicas = settings.Replicas
	}

	m := make(map[string]string, 0)
	if cfg.Envs != nil && len(cfg.Envs) > 0 {
		for _, e := range cfg.Envs {
			m[e.EnvKey] = e.EnvValue
		}
	}
	opt := &types.CreateDeploymentOpts{
		Envs:     m,
		Name:     appName,
		Replicas: int32(replicas),
		Tag:      rls.DockerUrl,
		IsLocal:  true,
		Memory:   settings.Plan.Info.Memory,
		Cpu:      settings.Plan.Info.Cpu,
	}
	if err := d.appRepo.UpdateEnvironmentVars(appName, m); err != nil {
		return nil, err
	}
	result, err := d.k8s.DeployService(opt)
	if err != nil {
		log.Println("failed to deploy service: ", err)
		return nil, ErrDeploymentFailed
	}
	if err := d.proxy.Set(appName, result.Address); err != nil {
		log.Println("failed to contact proxy server: ", err)
		return nil, ErrDeploymentFailed
	}
	result.Version = version
	result.Address = fmt.Sprintf("https://%s.%s", app.AppName, d.cfg.ServerUrl)
	return result, nil
}

func (d *defaultDeploymentRepository) HasRegistryAuthorization(req *types.AuthorizationRequest) ([]string, error) {
	return []string{
		"push", "pull",
	}, nil
}

func NewDeploymentRepository(db *gorm.DB, docker services.DockerService,
	k8s services.K8sService, pr proxy.Client,
	appRepo AppsRepository, storage services.StorageClient,
	cfg *config.Config) DeploymentRepository {
	return &defaultDeploymentRepository{
		db:         db,
		docker:     docker,
		k8s:        k8s,
		proxy:      pr,
		appRepo:    appRepo,
		storage:    storage,
		cfg:        cfg,
	}
}
