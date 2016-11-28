package config

//AuthConfig contains the information necessary for the API authentication code
// to initialize and function
type AuthConfig struct {
	Type   string `yaml:"type"`
	Config map[string]interface{}
}
