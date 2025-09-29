package configuration

import (
	"bufio"
	"os"
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type gherkinFormatterOptions struct {
	Intend bool
	Intendation int
}

// type Config struct {
// 	GherkinFormatter gherkinFormatter `toml:"gherkinFormatter"`
// }

func ReadConfiguration(path string) (gherkinFormatterOptions, error) {
	var cfg interface{}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	tomlBytes := make([]byte, 0)
	for scanner.Scan() {
		tomlBytes = append(tomlBytes, scanner.Bytes()...)
	}
	err = toml.Unmarshal(tomlBytes, &cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	return gherkinFormatterOptions{}, nil
}