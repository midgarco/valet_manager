package config

import (
	"os"
	"strings"

	"github.com/midgarco/env"
)

// Option allows you to provide addtional settings to override
type Option struct {
	Key   string
	Value string
}

// LoadEnv ...
func LoadEnv(path string, opts ...Option) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	err := env.Load(path + ".env")
	if err != nil {
		return err
	}
	// Set override options
	Override(opts...)
	return nil
}

// Get returns the value from the envronment config
func Get(key string) string {
	return os.Getenv(key)
}

// GetWithDefault returns the value from the environment config
// or returns a default value if the setting is empty
func GetWithDefault(key, def string) string {
	s := os.Getenv(key)
	if s == "" {
		s = def
	}
	return s
}

// Override environment config with additional options
func Override(opts ...Option) {
	for _, opt := range opts {
		os.Setenv(opt.Key, opt.Value)
	}
}

// GetBool returns a boolean configuration setting
func GetBool(key string) bool {
	v := GetWithDefault(key, "false")
	return (strings.ToLower(v) == "true" || v == "1")
}
