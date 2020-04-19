package types

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	AccountId           uint   `json:"account_id"`
	AppName             string `json:"app_name"`
	AccessUrl           string `json:"access_url"`
	RegistryDownloadUrl string `json:"registry_download_url"`
	LocalAccessUrl      string `json:"local_access_url"`
	GitUrl              string `json:"git_url"`

	Environments []Environment `json:"environments"`
	Account      *Account      `json:"account" gorm:"-" sql:"-"`
}

type Environment struct {
	gorm.Model
	AppId    uint   `json:"app_id"`
	ResId    uint   `json:"res_id"` // for resources
	EnvKey   string `json:"env_key"`
	EnvValue string `json:"env_value"`
}

type Release struct {
	gorm.Model
	AppId         uint   `json:"app_id"`
	VersionNumber int64  `json:"version"`
	DockerUrl     string `json:"docker_url"`
}

type ReleaseConfig struct {
	Envs    []Environment `json:"envs"`
	Version string
}

type Instance struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Name    string `json:"name"`
	Started string `json:"started"`
}

type Domain struct {
	gorm.Model
	AppId   uint   `json:"app_id"`
	Address string `json:"address"`
}

func NewDomain(appId uint, addr string) *Domain {
	return &Domain{AppId: appId, Address: addr}
}

type DeploymentSettings struct {
	gorm.Model
	AppId    uint  `json:"app_id"`
	Replicas uint  `json:"replicas"`
	PlanId   uint  `json:"plan_id"`
	Plan     *Plan `json:"plan" gorm:"-" sql:"-"`
}

func NewDeploymentSettings(appId, replicas uint) *DeploymentSettings {
	return &DeploymentSettings{AppId: appId, Replicas: replicas}
}

func NewRelease(appId uint, ref string, v int64) *Release {
	return &Release{
		AppId:         appId,
		DockerUrl: ref,
		VersionNumber: v,
	}
}


func NewEnvVariable(appId, resId uint, k, v string) *Environment {
	return &Environment{
		AppId:    appId,
		EnvKey:   k,
		EnvValue: v,
		ResId:    resId,
	}
}

func NewApp(name string, accountId uint) *App {
	return &App{
		AccountId: accountId,
		AppName:   name,
		AccessUrl: fmt.Sprintf("https://%s.hostgoapp.com", name),
		GitUrl:    fmt.Sprintf("https://git.hostgoapp.com/%s.git", name),
	}
}

func NewReleaseConfig(version string, envs []Environment) *ReleaseConfig {
	return &ReleaseConfig{
		Version: version, Envs: envs,
	}
}