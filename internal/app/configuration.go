package app

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Configuration represents the configuration yaml file
type Configuration struct {
	AuthEmail string `json:"auth_email"`
	AuthKey   string `json:"auth_key"`
	Zones     []Zone `json:"zones"`
}

// Zone represents a domain
type Zone struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	DNSRecords []DNSRecord `json:"dns_records"`
}

// DNSRecord represents a subdomain
type DNSRecord struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// ReadConfigurationFile parses a given file and returns a configuration object
func ReadConfigurationFile(filePath string) (*Configuration, error) {
	// Read configuration file
	configurationYaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal yaml
	var configuration Configuration
	if err = yaml.Unmarshal(configurationYaml, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

// WriteToFile writes a given configuration object to a specified file path
func (configuration *Configuration) WriteToFile(filePath string, fileMode os.FileMode) error {
	// Marshal yaml
	configurationYaml, err := yaml.Marshal(configuration)
	if err != nil {
		return err
	}

	// Write configuration file
	if err = ioutil.WriteFile(filePath, configurationYaml, fileMode); err != nil {
		return err
	}

	return nil
}
