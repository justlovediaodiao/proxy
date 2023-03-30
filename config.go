package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config define
type Config struct {
	Host          string   `json:"host"`
	Port          int      `json:"port"`
	Protocol      string   `json:"protocol"`
	PACHost       string   `json:"pac_host"`
	PACPort       int      `json:"pac_port"`
	ProxyCommands []string `json:"proxy_commands"`
	ProxyURL      string
}

func getProxyURL(protocol string, host string, port int) (string, error) {
	if protocol == "http" {
		return fmt.Sprintf("PROXY %s:%d;DIRECT", host, port), nil
	} else if protocol == "socks" {
		return fmt.Sprintf("SOCKS %s:%d;DIRECT", host, port), nil
	} else if protocol == "socks5" {
		return fmt.Sprintf("SOCKS5 %s:%d;DIRECT", host, port), nil
	}
	return "", fmt.Errorf("unspported proxy protocol %s", protocol)
}

// GetConfig read config from file.
func GetConfig() (*Config, error) {
	content, err := os.ReadFile("resources/config.json")
	if err != nil {
		return nil, err
	}
	var c = new(Config)
	if err = json.Unmarshal(content, c); err != nil {
		return nil, err
	}
	url, err := getProxyURL(c.Protocol, c.Host, c.Port)
	if err != nil {
		return nil, err
	}
	c.ProxyURL = url
	return c, nil
}
