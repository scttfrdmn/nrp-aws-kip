package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	AWS       AWSConfig       `yaml:"aws"`
	Node      NodeConfig      `yaml:"node"`
	Limits    LimitsConfig    `yaml:"limits"`
	Logging   LoggingConfig   `yaml:"logging"`
	Metrics   MetricsConfig   `yaml:"metrics"`
}

// AWSConfig contains AWS-specific configuration
type AWSConfig struct {
	Region            string            `yaml:"region"`
	VpcID             string            `yaml:"vpcId"`
	SubnetIDs         []string          `yaml:"subnetIds"`
	SecurityGroupIDs  []string          `yaml:"securityGroupIds"`
	KeyName           string            `yaml:"keyName"`
	IAMInstanceProfile string           `yaml:"iamInstanceProfile"`
	Tags              map[string]string `yaml:"tags"`
	UseSpotInstances  bool              `yaml:"useSpotInstances"`
	SpotMaxPrice      string            `yaml:"spotMaxPrice,omitempty"`
}

// NodeConfig contains virtual node configuration
type NodeConfig struct {
	Name               string            `yaml:"name"`
	OperatingSystem    string            `yaml:"operatingSystem"`
	CPU                string            `yaml:"cpu"`
	Memory             string            `yaml:"memory"`
	Pods               string            `yaml:"pods"`
	Labels             map[string]string `yaml:"labels"`
	Taints             []Taint           `yaml:"taints"`
}

// Taint represents a Kubernetes taint
type Taint struct {
	Key    string `yaml:"key"`
	Value  string `yaml:"value"`
	Effect string `yaml:"effect"`
}

// LimitsConfig contains resource limits and quotas
type LimitsConfig struct {
	MaxPods       int    `yaml:"maxPods"`
	MaxInstances  int    `yaml:"maxInstances"`
	CostLimit     string `yaml:"costLimit,omitempty"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Node.OperatingSystem == "" {
		cfg.Node.OperatingSystem = "Linux"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}
	if cfg.Metrics.Port == 0 {
		cfg.Metrics.Port = 10255
	}
	if cfg.Metrics.Path == "" {
		cfg.Metrics.Path = "/metrics"
	}

	return &cfg, nil
}
