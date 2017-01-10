package config

//APIConfig is a struct containing all the information necessary to set up the
// admin API for Portcullis
type APIConfig struct {
	Port        int        `yaml:"port"`
	Description string     `yaml:"description"`
	Auth        AuthConfig `yaml:"auth"`
	//TODO: TLS, etc
}
