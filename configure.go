package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
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

	// Select zones
	zones, err := getZones(authEmail, authKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("\nChoose which zones shall be used, if multiple zones are used separate them using commas:\n")

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
		log.Fatal("please select min. 1 zone")
	}

	// Choose DNS entries
	var selectedZones []Zone
	for _, i := range zoneSelectionIndexes {
		// Zone
		zone := zones[i]

		// Dns records
		dnsRecords, err := getDnsRecords(authEmail, authKey, zone)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("\nChoose dns records for " + zone.Name + ", if multiple dns records are used separate them using commas:\n")

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

		var selectedDnsRecords []DnsRecord
		for _, i := range dnsRecordSelectionIndexes {
			selectedDnsRecords = append(selectedDnsRecords, dnsRecords[i])
		}

		if len(selectedDnsRecords) > 0 {
			zone.DnsRecords = selectedDnsRecords
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

	configurationJson, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		log.Fatal("error")
	}

	// Write configuration file
	err = ioutil.WriteFile("configuration.json", configurationJson, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("\nConfiguration file has been created successfully.")
}

func getSelectedIndexes(selection string, length int) ([]int, error) {
	// Split into slice
	selectionStringSlice := strings.Split(selection, ",")
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
			} else {
				if itemInt > length {
					return nil, fmt.Errorf(item + " is not a valid selection")
				}

				// Check if slice already contains item
				if !contains(selectedIndexes, itemInt) {
					selectedIndexes = append(selectedIndexes, itemInt - 1)
				}
			}
		}
	}

	// Sort and return
	sort.Ints(selectedIndexes)
	return selectedIndexes, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getZones(authEmail string, authKey string) ([]Zone, error) {
	// Create request
	client := &http.Client{}
	uri := "https://api.cloudflare.com/client/v4/zones"
	req, _ := http.NewRequest("GET", uri, nil)

	// Set headers
	req.Header.Set("X-Auth-Email", authEmail)
	req.Header.Set("X-Auth-Key", authKey)
	req.Header.Set("Content-Type", "application/json")

	// Run request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check result
	if res.StatusCode == 200 {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal json
		var cloudflareResult CloudflareZones
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			return nil, err
		}

		// Get zones
		var zones []Zone
		for _, element := range cloudflareResult.Result {
			zones = append(zones, Zone{element.ID, element.Name, element.Status, nil})
		}

		return zones, nil
	} else {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal json
		var cloudflareResult CloudflareZones
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			return nil, err
		}

		// Check for errors
		if cloudflareResult.Success != true {
			for _, err := range cloudflareResult.Errors {
				return nil, fmt.Errorf(strconv.Itoa(err.Code) + " - " + err.Message)
			}
		}

		return nil, fmt.Errorf("unknown error")
	}
}

func getDnsRecords(authEmail string, authKey string, zone Zone) ([]DnsRecord, error) {
	// Create request
	client := &http.Client{}
	uri := "https://api.cloudflare.com/client/v4/zones/" + zone.ID + "/dns_records"
	req, _ := http.NewRequest("GET", uri, nil)

	// Set headers
	req.Header.Set("X-Auth-Email", authEmail)
	req.Header.Set("X-Auth-Key", authKey)
	req.Header.Set("Content-Type", "application/json")

	// Run request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check result
	if res.StatusCode == 200 {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal json
		var cloudflareResult CloudflareDnsRecords
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			return nil, err
		}

		// Get dns records
		var dnsRecords []DnsRecord
		for _, element := range cloudflareResult.Result {
			dnsRecords = append(dnsRecords, DnsRecord{element.ID, element.Type, element.Name, element.Proxied, element.TTL})
		}

		return dnsRecords, nil
	} else {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal json
		var cloudflareResult CloudflareDnsRecords
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			return nil, err
		}

		// Check for errors
		if cloudflareResult.Success != true {
			for _, err := range cloudflareResult.Errors {
				return nil, fmt.Errorf(strconv.Itoa(err.Code) + " - " + err.Message)
			}
		}

		return nil, fmt.Errorf("unknown error")
	}
}