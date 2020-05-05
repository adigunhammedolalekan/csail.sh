package minio

import "github.com/saas/hostgolang/pkg/res"

type minioResource struct {
	memory, cpu, storageSize float64
	envs                     map[string]string
}

func (m *minioResource) PublishableEnvs() map[string]string {
	return m.envs
}

func (m *minioResource) Name() string {
	return "minio"
}

func (m *minioResource) Envs() map[string]string {
	return m.envs
}

func (m *minioResource) Image() string {
	return "minio/minio"
}

func (m *minioResource) Quota() res.Quota {
	return res.Quota{Memory: m.memory, Cpu: m.cpu, StorageSize: m.storageSize}
}

func (m *minioResource) Port() int {
	return 9000
}

func (m *minioResource) Args() []string {
	return []string{"server", "/data"}
}

func (m *minioResource) URI() string {
	return ""
}


func Minio(mem, cpu, storageSize float64, envs map[string]string) res.Res {
	return &minioResource{memory: mem, cpu: cpu, storageSize: storageSize, envs: envs}
}
