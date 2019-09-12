package main

import (
	"strconv"
	"fmt"
	"os/exec"
	"strings"
)

func SetGlobal(p *Proxy) error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	err = reset(networks)
	if err != nil {
		return nil
	}
	for _, network := range networks {
		if p.Protocol == "http" {
			_, err = execute("networksetup", "-setwebproxy", network, p.Host, strconv.Itoa(p.Port))
		} else if p.Protocol == "socks" {
			_, err = execute("networksetup", "-setsocksfirewallproxy", network, p.Host, strconv.Itoa(p.Port))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func SetPAC(p *Proxy) error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	err = reset(networks)
	if err != nil {
		return nil
	}
	var url = fmt.Sprintf("http://%s:%d", p.PACHost, p.PACPort)
	for _, network := range networks {
		_, err = execute("networksetup", "-setautoproxyurl", network, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func Reset() error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	return reset(networks)
}

func reset(networks []string) error {
	var err error
	for _, network := range networks {
		_, err = execute("networksetup", "-setautoproxystate", network, "off")
		if err != nil {
			return err
		}
		_, err = execute("networksetup", "-setwebproxystate", network, "off")
		if err != nil {
			return err
		}
		_, err = execute("networksetup", "-setsecurewebproxystate", network, "off")
		if err != nil {
			return err
		}
		_, err = execute("networksetup", "-setsocksfirewallproxystate", network, "off")
		if err != nil {
			return err
		}
	}
	return nil
}

func listNetwork() ([]string, error) {
	result, err := execute("networksetup", "-listallnetworkservices")
	if err != nil {
		return nil, err
	}
	var networks = make([]string, 0, 2)
	for _, network := range strings.Split(string(result), "\n") {
		if network == "Wi-Fi" || network == "Ethernet" {
			networks = append(networks, network)
		}
	}
	return networks, nil
}

func execute(name string, args ...string) (string, error) {
	var cmd = exec.Command(name, args...)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(result), nil
}
