package res

type Res interface {
	Name() string
	Envs() map[string]string
	Image() string
	Quota() Quota
	PublishableEnvs() map[string]string
	Port() int
	Args() []string
}

type Quota struct {
	Memory, Cpu, StorageSize float64
}
