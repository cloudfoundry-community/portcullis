package config

//StoreConfig contains information required to connect to a database
type StoreConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}
