package repository

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	proxy "github.com/saas/hostgolang/pkg/proxyclient"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"io"
	"io/ioutil"
	"log"
)

var ErrCreateImage = errors.New("failed to create docker image for app. Please contact support")
var ErrPushImage = errors.New("failed to push image to remote repository. Please contact support")
var ErrDeploymentFailed = errors.New("deployment failed. Please contact support")
var ErrNoChangeToDeploy = errors.New("no changes to deploy")

type DeploymentRepository interface {
	CreateDeployment(app *types.App, reader io.Reader) (*types.DeploymentResult, error)
	UpdateEnvironmentVars(app *types.App, envs []types.Environment) error
	GetApplicationLogs(appName string) (string, error)
	CheckRelease(appId uint, r io.Reader) (string, error)
	GetRelease(appId uint) (*types.Release, error)
	CreateRelease(appId uint, checkSum string) error
}

type defaultDeploymentRepository struct {
	db *gorm.DB
	docker services.DockerService
	k8s services.K8sService
	proxy proxy.Client
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
	result, err := d.k8s.DeployService(dockerUrl, appName, map[string]string{}, true)
	if err != nil {
		log.Println("failed to deploy service: ", err)
		return nil, ErrDeploymentFailed
	}
	if err := d.proxy.Set(appName, result.Address); err != nil {
		log.Println("failed to contact proxy server: ", err)
		return nil, ErrDeploymentFailed
	}
	if err := d.CreateRelease(app.ID, checkSum); err != nil {
		log.Println("failed to create a new release: ", err)
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
	switch err {
	case nil:
		if release.LastCheckSum == checkSum {
			return checkSum, ErrNoChangeToDeploy
		}else {
			log.Println("Updating release from ", release.LastCheckSum, " to ", checkSum)
			release.LastCheckSum = checkSum
			return checkSum, d.updateRelease(release)
		}
	default:
		return checkSum, nil
	}
}

func (d *defaultDeploymentRepository) GetRelease(appId uint) (*types.Release, error) {
	r := &types.Release{}
	err := d.db.Table("releases").Where("app_id = ?", appId).First(r).Error
	if err != nil {
		return nil, err
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
	return d.db.Table("releases").Where("id = ?", r.ID).Update(r).Error
}

func (d *defaultDeploymentRepository) CreateRelease(appId uint, checkSum string) error {
	release, err := d.GetRelease(appId)
	if err == nil && release.LastCheckSum != "" {
		return err
	}
	release = types.NewRelease(appId, checkSum)
	return d.db.Create(release).Error
}

func NewDeploymentRepository(db *gorm.DB, docker services.DockerService, k8s services.K8sService, pr proxy.Client) DeploymentRepository {
	return &defaultDeploymentRepository{
		db:     db,
		docker: docker,
		k8s:    k8s,
		proxy: pr,
	}
}
