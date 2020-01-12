package config

type Config struct {
	ProxyServerAddress string `json:"proxy_server_address"`
	Registry RegistryConfig `json:"registry"`
}

type RegistryConfig struct {
	Url string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}
