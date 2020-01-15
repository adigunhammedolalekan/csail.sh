package repository

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/fn"
	"github.com/saas/hostgolang/pkg/res"
	"github.com/saas/hostgolang/pkg/res/minio"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"log"
	"strings"
)

var ErrProvisionResource = errors.New("failed to provision resource")

type ResourcesDeployment interface {
	DeployResource(opt *types.DeployResourcesOpt) (*types.ResourceDeploymentResult, error)
	GetResource(appId uint, resName string) (*types.Resource, error)
}

type defaultResourcesDeploymentRepo struct {
	db *gorm.DB
	appRepo AppsRepository
	accountRepo AccountRepository
	k8s services.ResourcesService
}

func (d *defaultResourcesDeploymentRepo) DeployResource(opt *types.DeployResourcesOpt) (*types.ResourceDeploymentResult, error) {
	app, err := d.appRepo.GetApp(opt.AppName)
	if err != nil {
		return nil, err
	}
	existing, err := d.GetResource(app.ID, opt.Name)
	if err == nil && existing.Name != "" {
		return nil, errors.New("resource has already been added to this application")
	}
	var r res.Res
	switch opt.Name {
	case "minio":
		key := fn.GenerateRandomString(64)
		accessKey, secretKey := key, d.reverseMd5(key)
		m := map[string]string{
			"MINIO_ACCESS_KEY": accessKey,
			"MINIO_SECRET_KEY": secretKey,
			"MINIO_PORT": "9000",
		}
		r = minio.Minio(opt.Memory, opt.Cpu, opt.StorageSize, m)
	default:
		return nil, errors.New("resources type not supported yet")
	}
	tx := d.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	resource := types.NewResource(r.Name(), app.ID)
	if err := tx.Create(resource).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, ErrProvisionResource
	}
	rCfg := types.NewResourceConfig(resource.ID, &types.Quota{
		Memory: opt.Memory, Cpu: opt.Cpu, StorageSize: opt.StorageSize,
	})
	if err := tx.Create(rCfg).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, ErrProvisionResource
	}
	envs := make([]types.ResourceEnv, 0)
	for k, v := range r.Envs() {
		rEnv := types.NewEnvVariable(app.ID, k, v)
		if err := tx.Create(rEnv).Error; err != nil {
			log.Println(err)
			tx.Rollback()
			return nil, ErrProvisionResource
		}
		envs = append(envs, *types.NewResourceEnv(resource.ID, rEnv.EnvKey, rEnv.EnvValue))
	}
	if err := tx.Commit().Error; err != nil {
		log.Println(err)
		return nil, ErrProvisionResource
	}
	result, err := d.k8s.DeployResource(app, envs, r)
	if err != nil {
		return nil, err
	}
	if err := d.updateHostEnvConfig(app.ID, fmt.Sprintf("%s_HOST", strings.ToUpper(r.Name())), result.Ip); err != nil {
		log.Println(err)
	}
	return result, nil
}

func (d *defaultResourcesDeploymentRepo) GetResource(appId uint, name string) (*types.Resource, error) {
	r := &types.Resource{}
	err := d.db.Table("resources").Where("name = ? AND app_id = ?", name, appId).First(r).Error
	if err != nil {
		return nil, ErrNotFound
	}
	return r, nil
}

func (d *defaultResourcesDeploymentRepo) updateHostEnvConfig(appId uint, key, hostIp string) error {
	env := types.NewEnvVariable(appId, key, hostIp)
	return d.db.Create(env).Error
}

func (d *defaultResourcesDeploymentRepo) reverseMd5(s string) string {
	reversed := reverse(s)
	m5 := md5.New()
	m5.Write([]byte(reversed))
	return fmt.Sprintf("%x", m5.Sum(nil))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func NewResourcesDeploymentRepository(db *gorm.DB, appRepo AppsRepository, accountRepo AccountRepository, k8s services.ResourcesService) ResourcesDeployment {
	return &defaultResourcesDeploymentRepo{
		db:          db,
		appRepo:     appRepo,
		accountRepo: accountRepo,
		k8s: k8s,
	}
}
