package proxy

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)


type Config struct {
	Proxy ProxyConfig `yaml:"proxy"`
}

type ProxyConfig struct {
	ListenPort int `yaml:"listen_port"`
	Backend BackendConfig `yaml:"backend"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

type BackendConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type TimeoutConfig struct {
	ReadTimeout time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

// Load Configuration from YAML
func LoadConfig (filepath string) (*Config, error) {

	configData, err := os.ReadFile(filepath) 
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration from file: %w", err)
	}

	// Parse configData to Config struct
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("Failed to load configuration from yaml file, check configuration file: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("Failed to Validate Proxy configuration, check the yaml file: %w", err)
	}

	return &config, nil
} 


func (c *Config) Validate() (error) {
	if c.Proxy.ListenPort <= 0 || c.Proxy.ListenPort >= 65535 {
		return fmt.Errorf("invalid listen_port: %d (must be 1-65535)", c.Proxy.ListenPort)
	}

	if c.Proxy.Backend.Host == "" {
		return fmt.Errorf("Ivalid Backend Host, it shouldnt be empty")
	}

	if c.Proxy.Backend.Port <= 0 || c.Proxy.Backend.Port >= 65535 {
		return fmt.Errorf("invalid Backend Port: %d (must be 1-65535)", c.Proxy.Backend.Port)
	}

	return nil
}

func (c *Config) GetBackendURL() string {
	host := c.Proxy.Backend.Host
	port := c.Proxy.Backend.Port
	return fmt.Sprintf("http://%s:%d", host, port)
}