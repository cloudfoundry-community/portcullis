package config

//DatabaseConfig contains information required to connect to a database
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Location string `yaml:"location"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
