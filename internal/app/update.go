package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudflare/cloudflare-go"
	externalip "github.com/glendc/go-external-ip"
)

// UpdateDNSRecords updates all dns records in a configuration (if necessary)
func UpdateDNSRecords(configuration Configuration, savedIPPath string, savedIPFileMode os.FileMode) error {
	// Get current ip address
	consensus := externalip.DefaultConsensus(nil, nil)
	currentIP, err := consensus.ExternalIP()
	if err != nil {
		return err
	}

	// Get stored ip address
	savedIP, err := ioutil.ReadFile(savedIPPath)
	if err != nil {
		savedIP = []byte("")
	}

	// IP address changed
	if string(savedIP) != currentIP.String() {
		// Update dns records
		fmt.Println("Updating dns records ...")
		for _, zone := range configuration.Zones {
			for _, record := range zone.DNSRecords {
				if err := setNewIP(configuration.AuthEmail, configuration.AuthKey, zone, record, currentIP.String()); err != nil {
					return err
				}
			}
		}

		// Write new ip address
		err = ioutil.WriteFile(savedIPPath, []byte(currentIP.String()), savedIPFileMode)
		if err != nil {
			return err
		}
	} else {
		// Message
		fmt.Println("Nothing to do here ...")
	}

	return nil
}

// setNewIP replaces the current ip with the one supplied as an argument
func setNewIP(authEmail string, authKey string, zone Zone, record DNSRecord, address string) error {
	// Cloudflare client
	api, err := cloudflare.New(authKey, authEmail)
	if err != nil {
		return err
	}

	// DNS record
	updateParams := cloudflare.UpdateDNSRecordParams{
		ID:      record.ID,
		Name:    record.Name,
		Type:    record.Type,
		Content: address,
	}

	// Update dns record
	if _, err = api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zone.ID), updateParams); err != nil {
		return err
	}

	fmt.Println("DNS record " + record.Name + " successfully updated")

	return nil
}
