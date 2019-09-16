package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

func SetGlobal(p *Proxy) error {
	return execute("resources/sysproxy.exe", "global", p.Host, strconv.Itoa(p.Port))
}

func SetPAC(p *Proxy) error {
	var addr = fmt.Sprintf("http://%s:%d/", p.Host, p.Port)
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
