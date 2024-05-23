package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/brharrelldev/weatherAPI/service"
	"gopkg.in/yaml.v3"
	"os"
)

func parseConfig(apiKey, path string) (*service.Config, error) {
	if path == "" {
		return nil, errors.New("path required")
	}

	if apiKey == "" {
		return nil, errors.New("api key required")
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("file could not be open due to %v", err)
	}

	defer f.Close()

	buf := make([]byte, 1024)

	offset, err := f.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %v", err)
	}

	content := buf[:offset]
	var conf *service.Config

	if err := yaml.NewDecoder(bytes.NewBuffer(content)).Decode(&conf); err != nil {
		return nil, fmt.Errorf("error decoding config %v", err)
	}

	return conf, nil
}
