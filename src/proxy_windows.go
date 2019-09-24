package main

import (
	"fmt"
	"os/exec"
)

func SetGlobal(p *Proxy) error {
	var addr string
	if p.Protocol == "http" {
		addr = fmt.Sprintf("%s:%d", p.Host, p.Port)
	} else if p.Protocol == "socks" {
		addr = fmt.Sprintf("socks=%s:%d", p.Host, p.Port)
	}
	return execute("resources/sysproxy.exe", "global", addr)
}

func SetPAC(p *Proxy) error {
	var addr = fmt.Sprintf("http://%s:%d/", p.PACHost, p.PACPort)
	return execute("resources/sysproxy.exe", "pac", addr)
}

func Reset() error {
	return execute("resources/sysproxy.exe", "set", "1", "-", "-", "-")
}

func execute(name string, args ...string) error {
	var cmd = exec.Command(name, args...)
	_, err := cmd.Output()
	return err
}
