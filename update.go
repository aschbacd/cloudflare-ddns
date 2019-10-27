package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	externalip "github.com/glendc/go-external-ip"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		fmt.Println("Cloudflare DDNS\n")
		fmt.Println("Updating dns records ...")

		// Update dns records
		for _, zone := range configuration.Zones {
			for _, record := range zone.DnsRecords {
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

func setNewIP(authEmail string, authKey string, zone Zone, record DnsRecord, address string) {
	// Record data
	var data = make(map[string]interface{})
	data["type"] = record.Type
	data["name"] = record.Name
	data["content"] = address
	data["ttl"] = record.TTL
	data["proxied"] = record.Proxied

	// Marshal data
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create request
	client := &http.Client{}
	uri := "https://api.cloudflare.com/client/v4/zones/" + zone.ID + "/dns_records/" + record.ID
	req, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(dataJson))

	// Set headers
	req.Header.Set("X-Auth-Email", authEmail)
	req.Header.Set("X-Auth-Key", authKey)
	req.Header.Set("Content-Type", "application/json")

	// Run request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Check result
	if res.StatusCode == 200 {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
		}

		// Unmarshal json
		var cloudflareResult CloudflareDnsRecord
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			log.Println(err.Error())
		}

		if cloudflareResult.Result.Content == address {
			fmt.Println("DNS record " + cloudflareResult.Result.Name + " successfully updated")
		} else {
			fmt.Println("DNS record did not update correctly")
		}

	} else {
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
		}

		// Unmarshal json
		var cloudflareResult CloudflareDnsRecord
		err = json.Unmarshal(body, &cloudflareResult)
		if err != nil {
			log.Println(err.Error())
		}

		// Check for errors
		if cloudflareResult.Success != true {
			for _, err := range cloudflareResult.Errors {
				log.Println(strconv.Itoa(err.Code) + " - " + err.Message)
			}
		}
	}
}