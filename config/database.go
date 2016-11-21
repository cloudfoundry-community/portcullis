package config

//DatabaseConfig contains information required to connect to a database
type DatabaseConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}
