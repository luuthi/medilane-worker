package config

type DBConfig struct {
	User     string `yaml:"DB_USER"`
	Password string `yaml:"DB_PASSWORD"`
	Driver   string `yaml:"DB_DRIVER"`
	Name     string `yaml:"DB_NAME"`
	Host     string `yaml:"DB_HOST"`
	Port     string `yaml:"DB_PORT"`
}
