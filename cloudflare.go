package main

type CloudflareZones struct {
	Result []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"result"`
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type CloudflareDnsRecords struct {
	Result []struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		Proxied bool   `json:"proxied"`
		TTL     int    `json:"ttl"`
	} `json:"result"`
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

type CloudflareDnsRecord struct {
	Result struct {
		ID         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxied    bool      `json:"proxied"`
		TTL        int       `json:"ttl"`
	} `json:"result"`
	Success  bool          `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}