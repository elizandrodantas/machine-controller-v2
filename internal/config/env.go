package config

import (
	"os"

	"github.com/spf13/viper"
)

type Env struct {
	DB_HOST         string `mapstructure:"DB_HOST"`
	DB_PORT         string `mapstructure:"DB_PORT"`
	DB_USER         string `mapstructure:"DB_USER"`
	DB_PASS         string `mapstructure:"DB_PASS"`
	DB_NAME         string `mapstructure:"DB_NAME"`
	DB_URL          string `mapstructure:"DB_URL"`
	WEB_PORT        string `mapstructure:"WEB_PORT"`
	TIMEOUT         string `mapstructure:"TIMEOUT"`
	ENVIRONMENT     string `mapstructure:"ENVIRONMENT"`
	JWT_TOKEN       string `mapstructure:"JWT_TOKEN"`
	JWT_EXPIRE_HOUR string `mapstructure:"JWT_EXPIRE_HOUR"`
	TOKEN_TYPE      string `mapstructure:"TOKEN_TYPE"`
}

// GetEnv retrieves the environment configuration from the .env file or environment variables.
// If the .env file exists, it reads the configuration from it using viper library.
// Otherwise, it retrieves the configuration from environment variables using viper.AutomaticEnv().
// The retrieved configuration is then unmarshaled into the 'env' struct.
// If any error occurs during the process, it is returned along with the 'env' struct.
func GetEnv() (env Env, err error) {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASS", "")
	viper.SetDefault("DB_NAME", "postgres")

	viper.SetDefault("WEB_PORT", ":3000")
	viper.SetDefault("TIMEOUT", "30")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("JWT_TOKEN", "default_token")
	viper.SetDefault("JWT_EXPIRE_HOUR", "24")
	viper.SetDefault("TOKEN_TYPE", "Bearer")
	viper.SetDefault("DB_URL", "")

	if _, err = os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.AddConfigPath(".")

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&env)
		return
	}

	viper.AutomaticEnv()
	err = viper.Unmarshal(&env)
	return
}
