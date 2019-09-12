package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Proxy struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	PACHost  string `json:"pac_host"`
	PACPort  int    `json:"pac_port"`
	Global   bool   `json:"global"`
}

func (p *Proxy) URL() string {
	if p.Protocol == "http" {
		return fmt.Sprintf("PROXY %s %d", p.Host, p.Port)
	} else if p.Protocol == "SOCKS" {
		return fmt.Sprintf("SOCKS %s %d", p.Host, p.Port)
	}
	panic(fmt.Sprintf("unspported proxy protocol %s", p.Protocol))
}

func GetConfig() (*Proxy, error) {
	content, err := ioutil.ReadFile("resources/config.json")
	if err != nil {
		return nil, err
	}
	var result = new(Proxy)
	if err = json.Unmarshal(content, result); err != nil {
		return nil, err
	}
	return result, nil
}

func SaveConfig(p *Proxy) error {
	content, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("resources/config.json", content, 0644)
}
