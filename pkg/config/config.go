package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	EXPOSEPORT     string `mapstructure:"EXPOSE_PORT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBName         string `mapstructure:"DB_NAME"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	JWT            string `mapstructure:"JWT_CODE"`
	AUTHTOKEN      string `mapstructure:"AUTH_TOKEN"`
	ACCOUNTSID     string `mapstructure:"ACCOUNT_SID"`
	SERVICESID     string `mapstructure:"SERVICE_SID"`
	RAZORPAYKEY    string `mapstructure:"RAZORPAY_KEY"`
	RAZORPAYSECRET string `mapstructure:"RAZORPAY_SECRET"`
}

var envs = []string{
	"EXPOSE_PORT",                                             //localport
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", //database
	"JWT_CODE",                                 //jwt
	"AUTH_TOKEN", "ACCOUNT_SID", "SERVICE_SID", //twilio details
	"RAZORPAY_KEY", "RAZORPAY_SECRET", //razorpay details
}

var config Config

func LoadConfig() (Config, error) {

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}

// to get the secred code for jwt
func GetJWTCofig() string {
	return config.JWT
}

func GetLocalPort() string {
	return config.EXPOSEPORT
}

func GetCofig() Config {
	return config
}
