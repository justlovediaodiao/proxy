package proxy

import (
	"fmt"
	"os"
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
	return execute("resources/sysproxy.exe", "global", addr, "<local>;192.168.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*")
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

// StartProxy start proxy process. donot block.
func StartProxy(cmd string) (*os.Process, error) {
	var arr = strings.Split(cmd, " ")
	var c = exec.Command(arr[0], arr[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return nil, err
	}
	return c.Process, nil
}

func execute(name string, args ...string) error {
	var cmd = exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
