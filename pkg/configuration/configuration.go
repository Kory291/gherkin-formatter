package configuration

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	IntendAnd bool
	Intendation int
}

var Configuration Config
var ConfigFileNotFoundError viper.ConfigFileNotFoundError

func setDefaults() {
	viper.SetDefault("intendation", 2)
	viper.SetDefault("intend-and", true)
}

func ReadConfiguration(path string) (*Config, error) {

	setDefaults()

	viper.SetConfigName("gherkinFormatter")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &ConfigFileNotFoundError) {
			fmt.Println("Config file not find - will use defauts")
		} else {
			return nil, errors.New("something else went wrong here")
		}
	}

	Configuration.IntendAnd = viper.GetBool("intend-and")
	Configuration.Intendation = viper.GetInt("intendation")

	return &Configuration, nil
}

func WriteConfiguration(path string) (error) {

	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	err := viper.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func SetConfiguration(key string, value any) {
	viper.Set(key, value)
}