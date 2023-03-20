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
	err := Reset()
	if err != nil {
		return nil
	}
	if c.Protocol == "http" {
		err = execute("gsettings", "set", "org.gnome.system.proxy.http", "host", c.Host)
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy.http", "port", strconv.Itoa(c.Port))
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy.https", "host", c.Host)
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy.https", "port", strconv.Itoa(c.Port))
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy", "mode", "manual")
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(c.Protocol, "socks") {
		err = execute("gsettings", "set", "org.gnome.system.proxy.socks", "host", c.Host)
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy.socks", "port", strconv.Itoa(c.Port))
		if err != nil {
			return err
		}
		err = execute("gsettings", "set", "org.gnome.system.proxy", "mode", "manual")
		if err != nil {
			return err
		}
	}
	return nil
}

// SetPAC set os proxy to pac.
func SetPAC(c *Config) error {
	err := Reset()
	if err != nil {
		return nil
	}
	var url = fmt.Sprintf("http://%s:%d", c.PACHost, c.PACPort)
	err = execute("gsettings", "set", "org.gnome.system.proxy", "autoconfig-url", url)
	if err != nil {
		return err
	}
	err = execute("gsettings", "set", "org.gnome.system.proxy", "mode", "auto")
	if err != nil {
		return err
	}
	return nil
}

// Reset clear os proxy settings.
func Reset() error {
	err := execute("gsettings", "set", "org.gnome.system.proxy", "mode", "none")
	return err
}

// StartProxy start proxy process. block until process exit.
func StartProxy(cmd string) error {
	var arr = strings.Split(cmd, " ")
	err := execute(arr[0], arr[1:]...)
	return err
}

func execute(name string, args ...string) error {
	var cmd = exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
