package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/justlovediaodiao/proxy"
)

func cmdGlobal(n int) {
	c, err := proxy.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = proxy.SetGlobal(c); err != nil {
		fmt.Println(err)
		return
	}

	defer proxy.Reset()
	fmt.Println("proxy set to global mode")

	var cmd string
	if n < len(c.ProxyCommands) {
		cmd = c.ProxyCommands[n]
	}
	if cmd != "" {
		fmt.Println("start proxy")
		process, err := proxy.StartProxy(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer process.Signal(os.Interrupt)
	}
	waitForExit()
}

func waitForExit() {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func cmdPAC(n int) {
	c, err := proxy.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = proxy.SetPAC(c); err != nil {
		fmt.Println(err)
		return
	}

	defer proxy.Reset()
	fmt.Println("proxy set to pac mode")

	var cmd string
	if n < len(c.ProxyCommands) {
		cmd = c.ProxyCommands[n]
	}
	if cmd != "" {
		fmt.Println("start proxy")
		process, err := proxy.StartProxy(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer process.Signal(os.Interrupt)
	}

	var addr = fmt.Sprintf("%s:%d", c.PACHost, c.PACPort)
	fmt.Println("start pac server")
	go func() {
		if err = proxy.StartServer(addr); err != nil {
			fmt.Println(err)
		}
	}()
	waitForExit()
}

func cmdOff() {
	if err := proxy.Reset(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("proxy cleared")
	}
}

func cmdUpdate() {
	c, err := proxy.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = proxy.UpdatePAC(c); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("pac updated")
	}
}

func cmdHelp() {
	fmt.Println("usage proxy [g[n]/global[n]/pac[n]/off/clear/update]")
}

func splitN(cmd string, prefix string) int {
	n := cmd[len(prefix):]
	if n == "" {
		return 0
	}
	if i, err := strconv.Atoi(n); err == nil && i >= 0 {
		return i
	}
	fmt.Printf("invalid arg '%s'", cmd)
	return -1
}

func main() {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var dir = filepath.Dir(path)
	err = os.Chdir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cmd string
	if len(os.Args) == 2 {
		cmd = os.Args[1]
	}
	switch {
	case strings.HasPrefix(cmd, "global"):
		if n := splitN(cmd, "global"); n != -1 {
			cmdGlobal(n)
		}
	case strings.HasPrefix(cmd, "g"):
		if n := splitN(cmd, "g"); n != -1 {
			cmdGlobal(n)
		}
	case strings.HasPrefix(cmd, "pac"):
		if n := splitN(cmd, "pac"); n != -1 {
			cmdPAC(n)
		}
	default:
		switch cmd {
		case "off", "clear":
			cmdOff()
		case "update":
			cmdUpdate()
		default:
			cmdHelp()
		}
	}
}
