package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cloudflare/cloudflare-go"
	externalip "github.com/glendc/go-external-ip"
)

func update(configuration Configuration) {
	// Get current ip address
	consensus := externalip.DefaultConsensus(nil, nil)
	currentIP, err := consensus.ExternalIP()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get stored ip address
	savedIP, err := ioutil.ReadFile("address.txt")
	if err != nil {
		savedIP = []byte("")
	}

	// IP address changed
	if string(savedIP) != currentIP.String() {
		// Message
		fmt.Println("Cloudflare DDNS")
		fmt.Println("Updating dns records ...")

		// Update dns records
		for _, zone := range configuration.Zones {
			for _, record := range zone.DNSRecords {
				setNewIP(configuration.AuthEmail, configuration.AuthKey, zone, record, currentIP.String())
			}
		}

		// Write new ip address
		err = ioutil.WriteFile("address.txt", []byte(currentIP.String()), 0644)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// Message
		fmt.Println("Nothing to do here ...")
	}
}

func setNewIP(authEmail string, authKey string, zone Zone, record DNSRecord, address string) {
	// Cloudflare client
	api, err := cloudflare.New(authKey, authEmail)
	if err != nil {
		log.Fatal("authentication failed")
	}

	// New dns record
	dnsRecord := cloudflare.DNSRecord{ID: record.ID, Type: record.Type, Name: record.Name, Proxied: record.Proxied, TTL: record.TTL, Content: address}

	// Update dns record
	err = api.UpdateDNSRecord(zone.ID, record.ID, dnsRecord)
	if err != nil {
		log.Fatal("could not update " + record.Name)
	}

	fmt.Println("DNS record " + record.Name + " successfully updated")
}
