package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

func configure() {
	// Banner
	println("Cloudflare DDNS - Configurator\n")

	// Authentication email
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter authentication email: ")
	authEmail, _ := reader.ReadString('\n')
	authEmail = strings.TrimRight(authEmail, "\r\n")

	if authEmail == "" {
		log.Fatal("authentication email cannot be empty")
	}

	// Authentication key
	fmt.Print("Enter authentication key: ")
	authKey, _ := reader.ReadString('\n')
	authKey = strings.TrimRight(authKey, "\r\n")

	if authKey == "" {
		log.Fatal("authentication key cannot be empty")
	}

	// Cloudflare client
	api, err := cloudflare.New(authKey, authEmail)
	if err != nil {
		log.Fatal("authentication failed")
	}

	zones, err := api.ListZones()
	if err != nil {
		log.Fatal("authentication failed")
	}

	fmt.Println("\nChoose which zones shall be used, if multiple zones are used separate them using commas:")

	for i, zone := range zones {
		fmt.Print("[" + strconv.Itoa(i+1) + "] " + zone.Name + "\n")
	}

	// User selection
	fmt.Print("\nSelection: ")
	zoneSelection, _ := reader.ReadString('\n')
	zoneSelection = strings.TrimRight(zoneSelection, "\r\n")

	// Get selected indexes
	zoneSelectionIndexes, err := getSelectedIndexes(zoneSelection, len(zones))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Check if items selected
	if len(zoneSelectionIndexes) < 0 {
		log.Fatal("select min. 1 zone")
	}

	// Choose DNS entries
	var selectedZones []Zone
	for _, i := range zoneSelectionIndexes {
		// Zone
		zone := Zone{ID: zones[i].ID, Name: zones[i].Name, Status: zones[i].Status}

		var dnsRecord cloudflare.DNSRecord

		// DNS records
		dnsRecords, err := api.DNSRecords(zones[i].ID, dnsRecord)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("\nChoose dns records for " + zones[i].Name + ", if multiple dns records are used separate them using commas:\n")

		for i, dnsRecord := range dnsRecords {
			fmt.Print("[" + strconv.Itoa(i+1) + "] " + dnsRecord.Name + "\n")
		}

		// User selection
		fmt.Print("\nSelection: ")
		dnsRecordSelection, _ := reader.ReadString('\n')
		dnsRecordSelection = strings.TrimRight(dnsRecordSelection, "\r\n")

		dnsRecordSelectionIndexes, err := getSelectedIndexes(dnsRecordSelection, len(dnsRecords))
		if err != nil {
			log.Fatal(err.Error())
		}

		// Add dns records
		var selectedDNSRecords []DNSRecord
		for _, i := range dnsRecordSelectionIndexes {
			selectedDNSRecords = append(selectedDNSRecords, DNSRecord{ID: dnsRecords[i].ID, Type: dnsRecords[i].Type, Name: dnsRecords[i].Name, Proxied: dnsRecords[i].Proxied, TTL: dnsRecords[i].TTL})
		}

		// Add zone to selected zones
		if len(selectedDNSRecords) > 0 {
			zone.DNSRecords = selectedDNSRecords
			selectedZones = append(selectedZones, zone)
		}
	}

	if len(selectedZones) == 0 {
		log.Fatal("no dns records have been selected")
	}

	// Create configuration file
	var configuration Configuration

	configuration.AuthEmail = authEmail
	configuration.AuthKey = authKey
	configuration.Zones = selectedZones

	configurationJSON, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		log.Fatal("error")
	}

	// Write configuration file
	err = ioutil.WriteFile("configuration.json", configurationJSON, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("\nConfiguration file has been created successfully.")
}

func getSelectedIndexes(selection string, length int) ([]int, error) {
	// Split into slice
	selectionStringSlice := strings.Split(selection, ",")
	var selectedItems []int
	var selectedIndexes []int

	// Check every item
	for _, item := range selectionStringSlice {
		item = strings.Replace(item, " ", "", -1)

		// Only add item if it has value
		if item != "" {
			// Convert to int
			itemInt, err := strconv.Atoi(item)
			if err != nil {
				return nil, fmt.Errorf(item + " is not a valid selection")
			}

			if itemInt > length {
				return nil, fmt.Errorf(item + " is not a valid selection")
			}

			// Check if slice already contains item
			if !contains(selectedItems, itemInt) {
				selectedItems = append(selectedItems, itemInt)
			}
		}
	}

	// Convert items to index (-1)
	for _, selectedItem := range selectedItems {
		selectedIndexes = append(selectedIndexes, selectedItem-1)
	}

	// Sort and return
	sort.Ints(selectedIndexes)
	return selectedIndexes, nil
}

// Checks if slice contains integer
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
