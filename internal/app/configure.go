package app

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aschbacd/cloudflare-ddns/internal/utils"
	"github.com/cloudflare/cloudflare-go"
)

// CreateConfiguration creates a configuration file including selected dns records
func CreateConfiguration(filePath string, fileMode os.FileMode) error {
	// Stdin reader
	reader := bufio.NewReader(os.Stdin)

	// Authentication email
	fmt.Print("Enter authentication email: ")
	authEmail, _ := reader.ReadString('\n')
	authEmail = strings.TrimRight(authEmail, "\r\n")

	if authEmail == "" {
		return fmt.Errorf("authentication email cannot be empty")
	}

	// Authentication key
	fmt.Print("Enter authentication key: ")
	authKey, _ := reader.ReadString('\n')
	authKey = strings.TrimRight(authKey, "\r\n")

	if authKey == "" {
		return fmt.Errorf("authentication key cannot be empty")
	}

	// Cloudflare client
	api, err := cloudflare.New(authKey, authEmail)
	if err != nil {
		return err
	}

	// List zones
	zones, err := api.ListZones()
	if err != nil {
		return err
	}

	if len(zones) > 0 {
		fmt.Println("\nChoose which zones shall be used, if multiple zones are used, separate them using commas:")
		for i, zone := range zones {
			fmt.Print("[" + strconv.Itoa(i+1) + "] " + zone.Name + "\n")
		}
	} else {
		return fmt.Errorf("no zones available for this account")
	}

	// User selection
	fmt.Print("\nSelection: ")
	zoneSelection, _ := reader.ReadString('\n')
	zoneSelection = strings.TrimRight(zoneSelection, "\r\n")

	// Get selected indexes
	zoneSelectionIndexes, err := getSelectedIndexes(zoneSelection, len(zones))
	if err != nil {
		return err
	}

	// Check if items selected
	if len(zoneSelectionIndexes) < 0 {
		return fmt.Errorf("min. 1 zone must be selected")
	}

	// DNS entries
	var selectedZones []Zone
	for _, i := range zoneSelectionIndexes {
		// Zone
		zone := Zone{ID: zones[i].ID, Name: zones[i].Name, Status: zones[i].Status}

		// DNS records
		dnsRecords, err := api.DNSRecords(zones[i].ID, cloudflare.DNSRecord{})
		if err != nil {
			return err
		}

		if len(dnsRecords) > 0 {
			fmt.Println("\nChoose dns records for " + zones[i].Name + ", if multiple dns records are used separate them using commas:\n")
			for i, dnsRecord := range dnsRecords {
				fmt.Print("[" + strconv.Itoa(i+1) + "] " + dnsRecord.Name + "\n")
			}
		} else {
			fmt.Println("no dns records for " + zone.Name)
		}

		// User selection
		fmt.Print("\nSelection: ")
		dnsRecordSelection, _ := reader.ReadString('\n')
		dnsRecordSelection = strings.TrimRight(dnsRecordSelection, "\r\n")

		dnsRecordSelectionIndexes, err := getSelectedIndexes(dnsRecordSelection, len(dnsRecords))
		if err != nil {
			return err
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
		return fmt.Errorf("no dns records have been selected")
	}

	// Create configuration object
	var configuration Configuration
	configuration.AuthEmail = authEmail
	configuration.AuthKey = authKey
	configuration.Zones = selectedZones

	// Write configuration file
	if err := configuration.WriteToFile(filePath, fileMode); err != nil {
		return err
	}

	return nil
}

// getSelectedIndexes returns a slice of integers by filtering a given user input string
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
			if !utils.ContainsInt(selectedItems, itemInt) {
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
