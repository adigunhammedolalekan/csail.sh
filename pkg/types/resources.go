package types

import "github.com/jinzhu/gorm"

type Resource struct {
	gorm.Model
	AppId uint `json:"app_id"`
	Name string `json:"name"`
}

func NewResource(name string, appId uint) *Resource {
	return &Resource{
		AppId: appId,
		Name:  name,
	}
}

type ResourceDeploymentResult struct {
	ID string `json:"id"`
	Ip string `json:"ip"`
	Port string `json:"port"`
}

type ResourceConfig struct {
	gorm.Model
	ResourceId uint `json:"resource_id"`
	Quota *Quota `json:"quota"`
}

func NewResourceConfig(resId uint, quota *Quota) *ResourceConfig {
	return &ResourceConfig{
		ResourceId: resId,
		Quota: quota,
	}
}

type ResourceEnv struct {
	gorm.Model
	ResourceId uint `json:"resource_id"`
	EnvKey string `json:"env_key"`
	EnvValue string `json:"env_value"`
}

func NewResourceEnv(resId uint, k, v string) *ResourceEnv {
	return &ResourceEnv{
		ResourceId: resId,
		EnvKey:     k,
		EnvValue:   v,
	}
}

type Quota struct {
	gorm.Model
	Memory float64 `json:"memory"`
	Cpu float64 `json:"cpu"`
	StorageSize float64 `json:"storage_size"`
}

func NewQuota(mem, cpu, ss float64) *Quota {
	return &Quota{
		Memory:      mem,
		Cpu:         cpu,
		StorageSize: ss,
	}
}