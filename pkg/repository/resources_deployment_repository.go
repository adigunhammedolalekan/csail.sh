package repository

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/saas/hostgolang/pkg/fn"
	"github.com/saas/hostgolang/pkg/res"
	"github.com/saas/hostgolang/pkg/res/minio"
	"github.com/saas/hostgolang/pkg/res/mongo"
	"github.com/saas/hostgolang/pkg/res/mysql"
	"github.com/saas/hostgolang/pkg/res/pg"
	"github.com/saas/hostgolang/pkg/services"
	"github.com/saas/hostgolang/pkg/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ErrProvisionResource = errors.New("failed to provision resource")
//go:generate mockgen -destination=../mocks/resource_deployment_repository_mock.go -package=mocks github.com/saas/hostgolang/pkg/repository ResourcesDeployment
type ResourcesDeployment interface {
	DeployResource(opt *types.DeployResourcesOpt) (*types.ResourceDeploymentResult, error)
	GetResource(appId uint, resName string) (*types.Resource, error)
	DeleteResource(app *types.App, resId uint, resName string) error
	DumpDatabase(app *types.App, resName string) (io.Reader, error)
	RestoreDatabase(app *types.App, resName string, data io.Reader) error
	GetResourceEnvs(appId, resId uint) ([]types.ResourceEnv, error)
}

type defaultResourcesDeploymentRepo struct {
	db          *gorm.DB
	appRepo     AppsRepository
	accountRepo AccountRepository
	k8s         services.ResourcesService
	docker services.DockerService
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
			"MINIO_PORT":       "9000",
		}
		r = minio.Minio(opt.Memory, opt.Cpu, opt.StorageSize, m)
	case "pg":
		key := fn.GenerateRandomString(35)
		username, password, dbName := key, d.reverseMd5(key), key[:20]
		dataPath := fmt.Sprintf("/var/lib/postgresl/data/pg-%s", app.AppName)
		m := map[string]string{
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       dbName,
			"PG_DATA":           dataPath,
		}
		r = pg.Postgres(opt.Memory, opt.Cpu, opt.StorageSize, m)
	case "mysql":
		key := fn.GenerateRandomString(50)
		username, password, dbName := fn.GenerateRandomString(30), d.reverseMd5(key), key[:20]
		m := map[string]string{
			"MYSQL_USER":          username,
			"MYSQL_PASSWORD":      password,
			"MYSQL_DATABASE":      dbName,
			"MYSQL_ROOT_PASSWORD": password,
		}
		r = mysql.NewMysql(opt.Memory, opt.Cpu, opt.StorageSize, m)
	case "mongo":
		key := fn.GenerateRandomString(50)
		username, password := key, d.reverseMd5(key)
		m := map[string]string{
			"MONGO_INITDB_ROOT_PASSWORD": password,
			"MONGO_INITDB_ROOT_USERNAME": username,
		}
		r = mongo.NewMongo(opt.Memory, opt.Cpu, opt.StorageSize, m)
	default:
		return nil, errors.New("resources type not supported yet")
	}
	tx := d.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	resource := types.NewResource(r.Name(), app.ID)
	if err := tx.Create(resource).Error; err != nil {
		tx.Rollback()
		return nil, ErrProvisionResource
	}
	rCfg := types.NewResourceConfig(resource.ID, &types.Quota{
		Memory: opt.Memory, Cpu: opt.Cpu, StorageSize: opt.StorageSize,
	})
	if err := tx.Create(rCfg).Error; err != nil {
		tx.Rollback()
		return nil, ErrProvisionResource
	}
	envs := make([]types.ResourceEnv, 0)
	for k, v := range r.Envs() {
		env := types.NewEnvVariable(app.ID, resource.ID, k, v)
		if err := tx.Create(env).Error; err != nil {
			tx.Rollback()
			return nil, ErrProvisionResource
		}
		resourceEnv := types.NewResourceEnv(resource.ID, env.EnvKey, env.EnvValue)
		if err := tx.Create(resourceEnv).Error; err != nil {
			tx.Rollback()
			return nil, ErrProvisionResource
		}
		envs = append(envs, *resourceEnv)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, ErrProvisionResource
	}
	appEnvs, _ := d.appRepo.GetEnvironmentVars(app.AppName)
	result, err := d.k8s.DeployResource(app, envs, appEnvs, r, true)
	if err != nil {
		return nil, err
	}
	if err := d.updateHostEnvConfig(app.ID, resource.ID, fmt.Sprintf("%s_HOST", strings.ToUpper(r.Name())), result.Ip); err != nil {
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

func (d *defaultResourcesDeploymentRepo) deleteResourceEnv(resId uint) error {
	if resId == 0 {
		return errors.New("resources does not exists")
	}
	return d.db.Table("environments").Where("res_id = ?", resId).Delete(&types.Environment{}).Error
}

func (d *defaultResourcesDeploymentRepo) DeleteResource(app *types.App, resId uint, resName string) error {
	if err := d.deleteResourceEnv(resId); err != nil {
		return err
	}
	if err := d.k8s.DeleteResource(app, resName); err != nil {
		return err
	}
	return d.db.Table("resources").Where("app_id = ? AND name = ?", app.ID, resName).Delete(&types.Resource{}).Error
}

func (d *defaultResourcesDeploymentRepo) DumpDatabase(app *types.App, resName string) (io.Reader, error) {
	r, err := d.GetResource(app.ID, resName)
	if err != nil {
		return nil, err
	}
	envs, err := d.GetResourceEnvs(app.ID, r.ID)
	if err != nil {
		return nil, err
	}
	switch resName {
	case "pg":
		cmds := make([]string, 0)
		cmds = append(cmds, "pg_dump", "-U")
		for _, e := range envs {
			if e.EnvKey == "POSTGRES_USER" {
				cmds = append(cmds, e.EnvValue)
			}
			if e.EnvKey == "POSTGRES_DB" {
				cmds = append(cmds, e.EnvValue)
			}
		}
		return d.k8s.Exec(app.AppName, "pg", cmds, nil)
	case "mysql":
		cmds := make([]string, 0)
		cmds = append(cmds, "mysqldump", "-u")
		databaseName := ""
		for _, e := range envs {
			if e.EnvKey == "MYSQL_USER" {
				cmds = append(cmds, e.EnvValue)
			}
			if e.EnvKey == "MYSQL_PASSWORD" {
				cmds = append(cmds, fmt.Sprintf("-p%s", e.EnvValue))
			}
			if e.EnvKey == "MYSQL_DATABASE" {
				databaseName = e.EnvValue
			}
		}
		cmds = append(cmds, databaseName)
		return d.k8s.Exec(app.AppName, "mysql", cmds, nil)
	default:
		return nil, errors.New("not yet implemented")
	}
}

func (d *defaultResourcesDeploymentRepo) RestoreDatabase(app *types.App, resName string, data io.Reader) error {
	r, err := d.GetResource(app.ID, resName)
	if err != nil {
		return err
	}
	envs, err := d.GetResourceEnvs(app.ID, r.ID)
	if err != nil {
		return err
	}
	switch resName {
	case "pg":
		filename, err := prepare("pg", data)
		if err != nil {
			return err
		}
		log.Println(filename)
		cmds := make([]string, 0)
		cmds = append(cmds,"-- psql", "-U")
		databaseName := ""
		for _, e := range envs {
			if e.EnvKey == "POSTGRES_USER" {
				cmds = append(cmds, e.EnvValue)
			}
			if e.EnvKey == "POSTGRES_DB" {
				databaseName = e.EnvValue
			}
		}
		cmds = append(cmds, "-d", databaseName, "<", "db-pg-tmp")
		r, err := d.k8s.Exec(app.AppName, "pg", cmds, data)
		if err != nil {
			return err
		}
		s, err := ioutil.ReadAll(r)
		log.Println(string(s), err)
		return nil
	default:
		return errors.New("not yet implemented")
	}
}

func prepare(resName string, data io.Reader) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	log.Println("current working dir ", wd)
	filename := filepath.Join(wd, fmt.Sprintf("db-%s-tmp", resName))
	fi, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fi, data); err != nil {
		return "", err
	}
	return filename, nil
}

func (d *defaultResourcesDeploymentRepo) GetResourceEnvs(appId, resId uint) ([]types.ResourceEnv, error) {
	data := make([]types.ResourceEnv, 0)
	err := d.db.Table("resource_envs").Where("resource_id = ?", resId).Find(&data).Error
	if err != nil {
		return nil, err
	}
	log.Println(data)
	return data, nil
}

func (d *defaultResourcesDeploymentRepo) updateHostEnvConfig(appId, resId uint, key, hostIp string) error {
	env := types.NewEnvVariable(appId, resId, key, hostIp)
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

func NewResourcesDeploymentRepository(db *gorm.DB, appRepo AppsRepository,
	accountRepo AccountRepository, k8s services.ResourcesService,
	docker services.DockerService) ResourcesDeployment {
	return &defaultResourcesDeploymentRepo{
		db:          db,
		appRepo:     appRepo,
		accountRepo: accountRepo,
		k8s:         k8s,
		docker: docker,
	}
}
