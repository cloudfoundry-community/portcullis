package config

//BrokerConfig contains configuration options used to set up a broker connection.
type BrokerConfig struct {
	Port         int    `yaml:"port"`
	CFAPIAddress string `yaml:"cf_api_address"`
	CFAdmin      string `yaml:"cf_admin"`
	CFPassword   string `yaml:"cf_password"`
}
