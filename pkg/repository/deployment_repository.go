package repository

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/config"
	proxy "github.com/saas/hostgolang/pkg/proxyclient"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ErrCreateImage = errors.New("failed to create docker image for app. Please contact support")
var ErrPushImage = errors.New("failed to push image to remote repository. Please contact support")
var ErrDeploymentFailed = errors.New("deployment failed. Please contact support")
var ErrNoChangeToDeploy = errors.New("no changes to deploy")
var ErrNotFound = errors.New("resources not found")
var ErrReleaseNotFound = errors.New("release not found")
var ErrNoDockerFile = errors.New("no Dockerfile found in this git repository")

const tempClonePath = "/tmp/git"

type DeploymentRepository interface {
	CreateDeployment(app *types.App, reader io.Reader) (*types.DeploymentResult, error)
	CreateGitDeployment(app *types.App, info *types.HookInfo) (*types.DeploymentResult, error)
	CreateDockerDeployment(app *types.App, dockerUrl string) (*types.DeploymentResult, error)
	UpdateEnvironmentVars(app *types.App, envs []types.Environment) error
	GetApplicationLogs(appName string) (string, error)
	CheckRelease(appId uint, r io.Reader) (string, error)
	GetRelease(appId uint) (*types.Release, error)
	CreateRelease(appId uint, checkSum string, data []byte) error
	CreateOrUpdateDeploymentSettings(appId, replicas uint) error
	GetDeploymentSettings(appId uint) (*types.DeploymentSettings, error)
	RollbackDeployment(appId uint, version string) (*types.DeploymentResult, error)
}

type defaultDeploymentRepository struct {
	db      *gorm.DB
	docker  services.DockerService
	k8s     services.K8sService
	proxy   proxy.Client
	appRepo AppsRepository
	storage services.StorageClient
	cfg *config.Config
	gitService services.GitService
}

func (d *defaultDeploymentRepository) CreateDeployment(app *types.App, reader io.Reader) (*types.DeploymentResult, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	builder := bytes.NewBuffer(data)

	checkSum, err := d.CheckRelease(app.ID, buf)
	if err == ErrNoChangeToDeploy {
		return nil, err
	}
	buildDir := fmt.Sprintf("%sBuild", app.AppName)
	appName := app.AppName

	dockerUrl, err := d.docker.BuildImage(context.Background(), buildDir, appName, builder)
	if err != nil {
		log.Println("failed to build docker image; ", err)
		return nil, ErrCreateImage
	}
	if err := d.docker.PushImage(context.Background(), dockerUrl); err != nil {
		log.Println("failed to push built image: ", err)
		return nil, ErrPushImage
	}
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
	if err := d.CreateRelease(app.ID, checkSum, data); err != nil {
		log.Println("failed to create a new release: ", err)
	}
	if err := d.CreateOrUpdateDeploymentSettings(app.ID, replicas); err != nil {
		log.Println("failed to create deployment settings: ", err)
	}
	rs, err := d.GetRelease(app.ID)
	if err == nil {
		result.Version = fmt.Sprintf("v%d", rs.VersionNumber)
	}
	return result, nil
}

func (d *defaultDeploymentRepository) CreateDockerDeployment(app *types.App, dockerUrl string) (*types.DeploymentResult, error) {
	rls, err := d.GetRelease(app.ID)
	if err == nil && rls.DockerUrl == dockerUrl {
		return nil, ErrNoChangeToDeploy
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
	rls.DockerUrl = dockerUrl
	rls.VersionNumber = rls.VersionNumber + 1
	if err := d.updateRelease(rls); err != nil {
		log.Println("failed to create a new release: ", err)
	}
	if err := d.CreateOrUpdateDeploymentSettings(app.ID, replicas); err != nil {
		log.Println("failed to create deployment settings: ", err)
	}
	rs, err := d.GetRelease(app.ID)
	if err == nil {
		result.Version = fmt.Sprintf("v%d", rs.VersionNumber)
	}
	return result, nil
}

func (d *defaultDeploymentRepository) UpdateEnvironmentVars(app *types.App, envs []types.Environment) error {
	return d.k8s.UpdateEnvs(app.AppName, envs)
}

func (d *defaultDeploymentRepository) GetApplicationLogs(appName string) (string, error) {
	return d.k8s.GetLogs(appName)
}

func (d *defaultDeploymentRepository) CheckRelease(appId uint, r io.Reader) (string, error) {
	checkSum := d.calculateCheckSum(r)
	release, err := d.GetRelease(appId)
	if err != nil || release.LastCheckSum == "" {
		return checkSum, nil
	}
	if release.LastCheckSum == checkSum {
		return checkSum, ErrNoChangeToDeploy
	}
	return checkSum, nil
}

func (d *defaultDeploymentRepository) GetRelease(appId uint) (*types.Release, error) {
	r := &types.Release{}
	err := d.db.Table("releases").Where("app_id = ?", appId).First(r).Error
	if err != nil {
		return r, ErrReleaseNotFound
	}
	return r, nil
}

func (d *defaultDeploymentRepository) calculateCheckSum(reader io.Reader) string {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (d *defaultDeploymentRepository) updateRelease(r *types.Release) error {
	newVersionNumber := r.VersionNumber + 1
	r.VersionNumber = newVersionNumber
	return d.db.Table("releases").Where("id = ?", r.ID).Update(r).Error
}

func (d *defaultDeploymentRepository) CreateRelease(appId uint, checkSum string, data []byte) error {
	release, err := d.GetRelease(appId)
	if err != nil || release.LastCheckSum == "" {
		// first time. create a new release
		rls := types.NewRelease(appId, checkSum, 1)
		return d.db.Create(rls).Error
	}
	release.LastCheckSum = checkSum
	if err := d.updateRelease(release); err != nil {
		return err
	}
	app, _ := d.appRepo.GetAppById(appId)
	envs, _ := d.appRepo.GetEnvironmentVars(app.AppName)
	cfg := &types.ReleaseConfig{Envs: envs, Version: fmt.Sprintf("%s:v%d", app.AppName, release.VersionNumber)}
	go func(storage services.StorageClient, reader []byte, r *types.ReleaseConfig) {
		if err := storage.Put(r.Version, reader); err != nil {
			log.Println("failed to store app release: ", err)
		}
		if err := storage.PutReleaseConfig(r.Version, r); err != nil {
			log.Println("failed to store app release configuration: ", err)
		}
	}(d.storage, data, cfg)
	return nil
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

	app, err := d.appRepo.GetAppById(appId)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s:%s", app.AppName, version)
	reader, err := d.storage.Get(key)
	if err != nil {
		return nil, err
	}
	cfg, err := d.storage.GetReleaseConfig(key)
	if err != nil {
		return nil, err
	}
	buildDir := fmt.Sprintf("%sBuild", app.AppName)
	appName := app.AppName

	dockerUrl, err := d.docker.BuildImage(context.Background(), buildDir, appName, reader)
	if err != nil {
		log.Println("failed to build docker image; ", err)
		return nil, ErrCreateImage
	}
	if err := d.docker.PushImage(context.Background(), dockerUrl); err != nil {
		log.Println("failed to push built image: ", err)
		return nil, ErrPushImage
	}
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
	result.Version = version
	return result, nil
}

func (d *defaultDeploymentRepository) CreateGitDeployment(app *types.App, info *types.HookInfo) (*types.DeploymentResult, error) {
	dir := filepath.Join(tempClonePath, info.RepoName)
	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}
	cloneUrl := fmt.Sprintf("%s/%s", d.cfg.GitServerUrl, info.RepoName)
	_, err := git.PlainClone(dir, true, &git.CloneOptions{
		URL:               cloneUrl,
		Auth:              &http.BasicAuth{Username: os.Getenv("MASTER_PASSWORD"), Password: os.Getenv("MASTER_PASSWORD")},
	})
	if err != nil {
		log.Println("failed to clone repo: ", err)
		return nil, ErrDeploymentFailed
	}
	d.gitService.WriteNotification(info.NotificationKey, fmt.Sprintf("building docker image for %s...", app.AppName))
	dockerUrl, err := d.docker.BuildImageFromGitRepository(context.Background(), dir, app.AppName)
	if err != nil {
		log.Println("failed to build image from git repo: ", err)
		return nil, ErrCreateImage
	}
	d.gitService.WriteNotification(info.NotificationKey, "pushing image to remote repository...")
	if err := d.docker.PushImage(context.Background(), dockerUrl); err != nil {
		log.Println("failed to push image: ", err)
		return nil, ErrPushImage
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
	settings.Plan = &types.Plan{
		Name:  "Test",
		Alias: "t1",
		Price: 0,
		Info:  types.PlanInfo{
			Memory: 0.1, Cpu: 0.1,
		},
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
	d.gitService.WriteNotification(info.NotificationKey, "creating deployment...")
	result, err := d.k8s.DeployService(opt)
	if err != nil {
		log.Println("failed to deploy service: ", err)
		return nil, ErrDeploymentFailed
	}
	if err := d.proxy.Set(appName, result.Address); err != nil {
		log.Println("failed to contact proxy server: ", err)
		return nil, ErrDeploymentFailed
	}
	d.gitService.WriteNotification(info.NotificationKey, "deployment updated!")
	d.gitService.WriteNotification(info.NotificationKey, "=====")
	d.gitService.WriteNotification(info.NotificationKey, fmt.Sprintf("https://%s.hostgoapp.com", app.AppName))
	d.gitService.WriteNotification(info.NotificationKey, "=====")
	d.gitService.WriteNotification(info.NotificationKey, "...Done")
	return result, nil
}

func NewDeploymentRepository(db *gorm.DB, docker services.DockerService,
	k8s services.K8sService, pr proxy.Client,
	appRepo AppsRepository, storage services.StorageClient,
	cfg *config.Config, gitService services.GitService) DeploymentRepository {
	return &defaultDeploymentRepository{
		db:      db,
		docker:  docker,
		k8s:     k8s,
		proxy:   pr,
		appRepo: appRepo,
		storage: storage,
		cfg: cfg,
		gitService: gitService,
	}
}
