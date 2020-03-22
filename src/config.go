package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Proxy struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	PACHost      string `json:"pac_host"`
	PACPort      int    `json:"pac_port"`
	ProxyCommand string `json:"proxy_command"`
}

func (p *Proxy) URL() string {
	if p.Protocol == "http" {
		return fmt.Sprintf("PROXY %s:%d;DIRECT", p.Host, p.Port)
	} else if p.Protocol == "socks" {
		return fmt.Sprintf("SOCKS %s:%d;DIRECT", p.Host, p.Port)
	} else if p.Protocol == "socks5" {
		return fmt.Sprintf("SOCKS5 %s:%d;DIRECT", p.Host, p.Port)
	} else if p.Protocol == "socks4" {
		return fmt.Sprintf("SOCKS4 %s:%d;DIRECT", p.Host, p.Port)
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
