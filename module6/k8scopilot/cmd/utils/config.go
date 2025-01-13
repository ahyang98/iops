package utils

import "fmt"

type Config struct {
	store map[string]any
}

func NewConfig() *Config {
	return &Config{store: make(map[string]any)}
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}

func (c *Config) Set(key string, value any) {
	c.store[key] = value
}

func (c *Config) Get(key string) (any, error) {
	v, ok := c.store[key]
	if !ok {
		return nil, fmt.Errorf("key %s not exist", key)
	}
	return v, nil
}
