package configmanager

import (
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigValidate struct {
	TitleMax       int `yaml:"title_max"`
	DescriptionMax int `yaml:"description_max"`
	PhotosMax      int `yaml:"photos_max"`
}

type Config struct {
	BindAddr string `yaml:"bind_addr"`
	DB       struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
	Validate struct {
		TitleMax       int `yaml:"title_max"`
		DescriptionMax int `yaml:"description_max"`
		PhotosMax      int `yaml:"photos_max"`
	} `yaml:"validate"`
	Paginate struct {
		AdvertsPageSize int `yaml:"adverts_page_size"`
	} `yaml:"paginate"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}

func (cm *Config) Init(configPath string) error {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cm)
	return err
}
