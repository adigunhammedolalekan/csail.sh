package types

type NewAccountOpts struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	CompanyName    string `json:"company_name"`
	CompanyWebsite string `json:"company_website"`
}

type AuthenticateAccountOpts struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DeploymentResult struct {
	Address string `json:"address"`
	Version string `json:"version"`
}

type CreateDeploymentOpts struct {
	Envs        map[string]string
	Name        string
	Replicas    int32
	Tag         string
	IsLocal     bool
	Memory, Cpu float64
}

type DeployResourcesOpt struct {
	AppName     string  `json:"app_name"`
	Name        string  `json:"name"`
	Memory      float64 `json:"memory"`
	Cpu         float64 `json:"cpu"`
	StorageSize float64 `json:"storage_size"`
}
