package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/client/notify"
	"reflect"
	"strings"
)

const (
	// DefaultConfigurationName is the default name of configuration
	defaultConfigurationName = "gocron"

	// DefaultConfigurationPath the default location of the configuration file
	defaultConfigurationPath = "/etc/gocron"
)

// Config defines everything needed for server to deal with external services
type Config struct {
	DatabaseOptions *database.Options
	NotifyOptions   *notify.Options
}

// newConfig creates a default non-empty Config
func New() *Config {
	return &Config{
		DatabaseOptions: database.NewDatabaseOptions(),
		NotifyOptions:   notify.NewNotifyOptions(),
	}
}

// TryLoadFromDisk loads configuration from default location after server startup
// return nil error if configuration file not exists
func TryLoadFromDisk() (*Config, error) {
	viper.SetConfigName(defaultConfigurationName)
	viper.AddConfigPath(defaultConfigurationPath)

	// Load from current working directory, only used for debugging
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// convertToMap simply converts config to map[string]bool
// to hide sensitive information
func (conf *Config) ToMap() map[string]bool {
	conf.stripEmptyOptions()
	result := make(map[string]bool, 0)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}

// Remove invalid options before serializing to json or yaml
func (conf *Config) stripEmptyOptions() {
	if conf.DatabaseOptions != nil && conf.DatabaseOptions.Host == "" {
		conf.DatabaseOptions = nil
	}
	if conf.NotifyOptions != nil && (conf.NotifyOptions.WebhookOpt != nil || conf.NotifyOptions.MailOptions != nil) {
		conf.NotifyOptions = nil
	}
}
