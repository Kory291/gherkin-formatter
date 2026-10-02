package configuration

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	IndentAnd   bool
	Indentation int
	SortTags    bool
}

var Configuration Config
var ConfigFileNotFoundError viper.ConfigFileNotFoundError

func setDefaults() {
	viper.SetDefault("indentation", 2)
	viper.SetDefault("indent-and", true)
	viper.SetDefault("sort-tags", true)
}

func setConfigPaths(path string) {
	viper.SetConfigName("gherkinFormatter")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
}

func ReadConfiguration(path string) (*Config, error) {

	setDefaults()
	setConfigPaths(path)

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &ConfigFileNotFoundError) {
			fmt.Println("Config file not find - will use defauts")
		} else {
			return nil, errors.New("something else went wrong here")
		}
	}

	Configuration.IndentAnd = viper.GetBool("indent-and")
	Configuration.Indentation = viper.GetInt("indentation")
	Configuration.SortTags = viper.GetBool("sort-tags")

	return &Configuration, nil
}

func WriteConfiguration(path string) error {

	setConfigPaths(path)

	err := viper.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func SetConfiguration(key string, value any) {
	viper.Set(key, value)
}

func CreateConfiguration(configDir string) error {
	setDefaults()
	viper.SetConfigType("toml")
	return WriteConfiguration(configDir)
}

func PrintConfiguration(configuration *Config) {
	fmt.Println("Configuration read:")
	fmt.Printf("indent-and:\t%t\n", configuration.IndentAnd)
	fmt.Printf("indentation:\t%d\n", configuration.Indentation)
	fmt.Printf("sort-tags:\t%t\n", configuration.SortTags)
}
