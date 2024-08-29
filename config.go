package config

import (
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/jinzhu/copier"
	log "github.com/rowdyroad/go-simple-logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LoadConfigFromFile loading config from yaml file
func LoadConfigFromFile(config interface{}, configFileName string, defaultValue interface{}) (string, error) {
	log.Debugf("Reading configuration from '%s'", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Warn("Configuration not found")
		if defaultValue != nil {
			log.Warn("Default value is defined. Using it.")
			copier.Copy(config, defaultValue)
			return "", nil
		}
		return "", fmt.Errorf("os.Open(configFileName): %v", err)
	}

	if err := LoadConfigFromReader(config, configFile, defaultValue); err != nil {
		return "", fmt.Errorf("LoadConfigFromReader(config, configFile, defaultValue): %v", err)
	}

	customConfigFileName := filepath.Join(
		filepath.Dir(configFileName),
		strings.TrimSuffix(filepath.Base(configFileName), filepath.Ext(configFileName))+".custom"+filepath.Ext(configFileName),
	)
	log.Debugf("Try to read custom configuration from '%s'...", customConfigFileName)
	customConfigFile, err := os.Open(customConfigFileName)
	if err == nil {
		log.Debugf("Reading custom configuration from '%s'", customConfigFileName)
		if err = LoadConfigFromReader(config, customConfigFile, defaultValue); err != nil {
			return "", fmt.Errorf("LoadConfigFromReader(config, customConfigFile, defaultValue): %v", err)
		}
		log.Debug("Config loaded successfully with custom config file")
		return customConfigFileName, nil
	}

	log.Debug("Config loaded successfully")
	return configFileName, nil
}

func LoadConfigFromReader(config interface{}, configReader io.Reader, defaultValue interface{}) error {
	if err := yaml.NewDecoder(configReader).Decode(config); err != nil {
		log.Warn("Configuration incorrect ")
		if defaultValue != nil {
			log.Warn("Default value is defined. Use it.")
			copier.Copy(config, defaultValue)
			return nil
		}
		return fmt.Errorf("yaml.NewDecoder(configFile).Decode(config): %v", err)
	}

	return nil
}

// LoadConfig from command line argument
func LoadConfig(config interface{}, defaultFilename string, defaultValue interface{}) (string, error) {
	var configFile string
	flag.StringVar(&configFile, "c", defaultFilename, "Config file")
	flag.StringVar(&configFile, "config", defaultFilename, "Config file")
	flag.Parse()
	return LoadConfigFromFile(config, configFile, defaultValue)
}
