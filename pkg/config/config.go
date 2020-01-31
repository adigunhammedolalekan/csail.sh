package config

type Config struct {
	ProxyServerAddress string         `json:"proxy_server_address"`
	Registry           RegistryConfig `json:"registry"`
	GitServerUrl string `json:"git_server_url"`
	GitTcpAddr string `json:"git_tcp_port"`
}

type RegistryConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}
