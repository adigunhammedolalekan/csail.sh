package types

type NewAccountOpts struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	CompanyName    string `json:"company_name"`
	CompanyWebsite string `json:"company_website"`
}

type AuthenticateAccountOpts struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type DeploymentResult struct {
	Address string `json:"address"`
}