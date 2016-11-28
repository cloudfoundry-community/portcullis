package config

//APIConfig is a struct containing all the information necessary to set up the
// admin API for Portcullis
type APIConfig struct {
	Auth AuthConfig `yaml:"auth"`
	//TODO: TLS, etc
}
