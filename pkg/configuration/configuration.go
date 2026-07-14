package configuration

type ThingiverseConfiguration struct {
	ClientId     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	AccessToken  string `json:"access_token" yaml:"access_token"`
}
