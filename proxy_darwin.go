package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// SetGlobal set os proxy to http or socks.
func SetGlobal(c *Config) error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	err = reset(networks)
	if err != nil {
		return nil
	}
	for _, network := range networks {
		if c.Protocol == "http" {
			err = execute("networksetup", "-setwebproxy", network, c.Host, strconv.Itoa(c.Port))
			if err != nil {
				return err
			}
			err = execute("networksetup", "-setsecurewebproxy", network, c.Host, strconv.Itoa(c.Port))
			if err != nil {
				return err
			}

		} else if strings.HasPrefix(c.Protocol, "socks") {
			err = execute("networksetup", "-setsocksfirewallproxy", network, c.Host, strconv.Itoa(c.Port))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetPAC set os proxy to pac.
func SetPAC(c *Config) error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	err = reset(networks)
	if err != nil {
		return nil
	}
	var url = fmt.Sprintf("http://%s:%d", c.PACHost, c.PACPort)
	for _, network := range networks {
		err = execute("networksetup", "-setautoproxyurl", network, url)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reset clear os proxy settings.
func Reset() error {
	networks, err := listNetwork()
	if err != nil {
		return err
	}
	return reset(networks)
}

// StartProxy start proxy process. block until process exit.
func StartProxy(cmd string) error {
	var arr = strings.Split(cmd, " ")
	err := execute(arr[0], arr[1:]...)
	return err
}

func reset(networks []string) error {
	var err error
	for _, network := range networks {
		err = execute("networksetup", "-setautoproxystate", network, "off")
		if err != nil {
			return err
		}
		err = execute("networksetup", "-setwebproxystate", network, "off")
		if err != nil {
			return err
		}
		err = execute("networksetup", "-setsecurewebproxystate", network, "off")
		if err != nil {
			return err
		}
		err = execute("networksetup", "-setsocksfirewallproxystate", network, "off")
		if err != nil {
			return err
		}
	}
	return nil
}

func listNetwork() ([]string, error) {
	result, err := exec.Command("networksetup", "-listallnetworkservices").Output()
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

func execute(name string, args ...string) error {
	var cmd = exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
