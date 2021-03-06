package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetGlobal set os proxy to http or socks.
func SetGlobal(c *Config) error {
	var addr string
	if c.Protocol == "http" {
		addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	} else if strings.HasPrefix(c.Protocol, "socks") {
		addr = fmt.Sprintf("socks=%s:%d", c.Host, c.Port)
	}
	return execute("resources/sysproxy.exe", "global", addr)
}

// SetPAC set os proxy to pac.
func SetPAC(c *Config) error {
	var addr = fmt.Sprintf("http://%s:%d/", c.PACHost, c.PACPort)
	return execute("resources/sysproxy.exe", "pac", addr)
}

// Reset clear os proxy settings.
func Reset() error {
	return execute("resources/sysproxy.exe", "set", "1", "-", "-", "-")
}

// StartProxy start proxy process. block until process exit.
func StartProxy(cmd string) error {
	var arr = strings.Split(cmd, " ")
	return execute(arr[0], arr[1:]...)
}

func execute(name string, args ...string) error {
	var cmd = exec.Command(name, args...)
	_, err := cmd.Output()
	return err
}
