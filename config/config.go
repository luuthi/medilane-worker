package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"medilane-worker/logger"
	"os"
)

type Config struct {
	DB         DBConfig             `yaml:"DATABASE"`
	Logger     logger.ConfigLogging `yaml:"LOGGER"`
	REDIS      Redis                `yaml:"REDIS"`
	FcmKeyPath string               `yaml:"FCM_KEY"`
}

type Redis struct {
	URL      string `json:"URL" yaml:"URL"`
	DB       int    `json:"DB" yaml:"DB"`
	Password string `json:"PASSWORD" yaml:"PASSWORD"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = "/app/config.yaml"
	}
	//err := godotenv.Load(configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Println("Error loading yaml file")
	}

	return &conf
}
