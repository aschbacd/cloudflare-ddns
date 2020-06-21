package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Configuration represents the configuration json file
type Configuration struct {
	AuthEmail string `json:"auth_email"`
	AuthKey   string `json:"auth_key"`
	Zones     []Zone `json:"zones"`
}

// Zone represents a domain
type Zone struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Status     string      `json:"status"`
	DNSRecords []DNSRecord `json:"dns_records"`
}

// DNSRecord represents a subdomain
type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

// ReadConfigurationFile parses a given file and returns a configuration object
func ReadConfigurationFile(filePath string) (*Configuration, error) {
	// Read configuration file
	configurationJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal json
	var configuration Configuration
	if err = json.Unmarshal(configurationJSON, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

// WriteToFile writes a given configuration object to a specified file path
func (configuration *Configuration) WriteToFile(filePath string, fileMode os.FileMode) error {
	// Marshal json
	configurationJSON, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		return err
	}

	// Write configuration file
	if err = ioutil.WriteFile(filePath, configurationJSON, fileMode); err != nil {
		return err
	}

	return nil
}
