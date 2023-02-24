package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Metrics []Metric `yaml:"metrics"`
}

type Metric struct {
	Name       string `yaml:"name"`
	Help       string
	MetricType string `yaml:"type"`
	Labels     []string
	Generators []Generator
}
type Generator struct {
	Value  float64
	Freq   int64
	Labels map[string]string
	Method string
}

func ParseMetrics() config {
	var config config
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("configFile read err %v ", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}
	return config

}
