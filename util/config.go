package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Stripe   StripeConfig
	Security SecurityConfig
}
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type DatabaseConfig struct {
	DBDriver   string `yaml:"driver"`
	DBUsername string `yaml:"DBUsername"`
	DBName     string `yaml:"DBName"`
	DBPassword string `yaml:"DBPassword"`
	DBSchema   string `yaml:"DBSchema"`
}
type StripeConfig struct {
	StripeAPI      string `yaml:"StripeAPI"`
	CancelURL      string `yaml:"CancelURL"`
	SuccessURL     string `yaml:"SuccessURL"`
	EndpointSecret string `yaml:"EndpointSecret"`
}
type SecurityConfig struct {
	OAuth2 OAuth2
}
type OAuth2 struct {
	ClientID       string   `yaml:"ClientID"`
	ApiRedirectURL string   `yaml:"ApiRedirectUrl"`
	WebRedirectURL string   `yaml:"WebRedirectUrl"`
	IssuerUrl      string   `yaml:"IssuerUrl"`
	Scopes         []string `yaml:"Scopes"`
	UrlBase        string   `yaml:"UrlBase"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("subscription_manager")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
