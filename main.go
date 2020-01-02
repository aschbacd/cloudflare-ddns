package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

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

// Configuration represents the configuration json file
type Configuration struct {
	AuthEmail string `json:"auth_email"`
	AuthKey   string `json:"auth_key"`
	Zones     []Zone `json:"zones"`
}

func main() {
	// Check arguments
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		// Load configuration
		configurationJSON, err := ioutil.ReadFile("configuration.json")
		if err != nil {
			fmt.Println("Cloudflare DDNS")
			fmt.Println("no configuration file available")
			fmt.Println("run \"" + os.Args[0] + " --configure\" to create a configuration file")
			os.Exit(1)
		}

		// Unmarshal json
		var configuration Configuration
		err = json.Unmarshal(configurationJSON, &configuration)
		if err != nil {
			fmt.Println("please check configuration file syntax")
			log.Fatal(err.Error())
		}

		// Update dns entries
		update(configuration)
	} else {
		// Process arguments
		switch arguments[0] {
		case "--configure":
			// Start configurator
			configure()
		}
	}
}
