package config

//BrokerConfig contains configuration options used to set up a broker connection.
type BrokerConfig struct {
	Port int `yaml:"port"`
}
