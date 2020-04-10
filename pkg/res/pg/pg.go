package pg

import "github.com/saas/hostgolang/pkg/res"

type pgResource struct {
	memory, cpu, storageSize float64
	envs                     map[string]string
}

func (p *pgResource) Name() string {
	return "pg"
}

func (p *pgResource) Envs() map[string]string {
	return p.envs
}

func (p *pgResource) Image() string {
	return "postgres:9.6"
}

func (p *pgResource) Quota() res.Quota {
	return res.Quota{
		Memory:      p.memory,
		Cpu:         p.cpu,
		StorageSize: p.storageSize,
	}
}

func (p *pgResource) PublishableEnvs() map[string]string {
	return p.envs
}

func (p *pgResource) Port() int {
	return 5432
}

func (p *pgResource) Args() []string {
	return []string{}
}

func Postgres(mem, cpu, ss float64, envs map[string]string) res.Res {
	return &pgResource{
		memory: mem, cpu: cpu, storageSize: ss, envs: envs,
	}
}
