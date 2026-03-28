package proxy

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Proxy        ProxyConfig        `yaml:"proxy"`
	ControlPlane ControlPlaneConfig `yaml:"control_plane"`
}

type ProxyConfig struct {
	ID            string        `yaml:"id"`
	Version       string        `yaml:"version"`
	ListenPort    int           `yaml:"listen_port"`
	AdvertiseAddr string        `yaml:"advertise_addr"`
	Backend       BackendConfig `yaml:"backend"`
	Timeout       TimeoutConfig `yaml:"timeout"`
}

type ControlPlaneConfig struct {
	Address string `yaml:"address"`
}

type BackendConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TimeoutConfig struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Load Configuration from YAML
func LoadConfig(filepath string) (*Config, error) {

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

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Proxy.ID) == "" {
		return fmt.Errorf("invalid proxy id: it should not be empty")
	}

	if strings.TrimSpace(c.Proxy.Version) == "" {
		return fmt.Errorf("invalid proxy version: it should not be empty")
	}

	if err := validatePort(c.Proxy.ListenPort, "proxy listen_port"); err != nil {
		return err
	}

	if err := validateAddress(c.Proxy.AdvertiseAddr, "proxy advertise_addr"); err != nil {
		return err
	}

	if strings.TrimSpace(c.Proxy.Backend.Host) == "" {
		return fmt.Errorf("invalid backend host: it should not be empty")
	}

	if err := validatePort(c.Proxy.Backend.Port, "backend port"); err != nil {
		return err
	}

	if err := validateAddress(c.ControlPlane.Address, "control_plane.address"); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetBackendURL() string {
	host := c.Proxy.Backend.Host
	port := c.Proxy.Backend.Port
	return fmt.Sprintf("http://%s:%d", host, port)
}

func validatePort(port int, fieldName string) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid %s: %d (must be 1-65535)", fieldName, port)
	}

	return nil
}

func validateAddress(address string, fieldName string) error {
	if strings.TrimSpace(address) == "" {
		return fmt.Errorf("invalid %s: it should not be empty", fieldName)
	}

	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid %s %q: expected host:port: %w", fieldName, address, err)
	}

	if strings.TrimSpace(host) == "" {
		return fmt.Errorf("invalid %s %q: host should not be empty", fieldName, address)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid %s %q: port must be numeric", fieldName, address)
	}

	return validatePort(port, fieldName+" port")
}
