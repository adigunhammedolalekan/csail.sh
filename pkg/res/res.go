package res

type Res interface {
	Name() string
	Envs() map[string]string
	Image() string
	Quota() Quota
	PublishableEnvs() map[string]string
	Port() int
	Args() []string
	URI() string
}

type Quota struct {
	Memory, Cpu, StorageSize float64
}

type genericResource struct {
	port                     int
	memory, cpu, storageSize float64
	envs                     map[string]string
	name, image              string
}

func NewGenericResource(mem, cpu, ss float64, envs map[string]string, port int, name, image string) Res {
	return &genericResource{
		port:        port,
		memory:      mem,
		cpu:         cpu,
		storageSize: ss,
		envs:        envs,
		name:        name, image: image,
	}
}

func (p *genericResource) Name() string {
	return p.name
}

func (p *genericResource) Envs() map[string]string {
	return p.envs
}

func (p *genericResource) Image() string {
	return p.image
}

func (p *genericResource) Quota() Quota {
	return Quota{
		Memory:      p.memory,
		Cpu:         p.cpu,
		StorageSize: p.storageSize,
	}
}

func (p *genericResource) PublishableEnvs() map[string]string {
	return p.envs
}

func (p *genericResource) Port() int {
	return p.port
}

func (p *genericResource) Args() []string {
	return []string{}
}

func (p *genericResource) URI() string {
	return ""
}