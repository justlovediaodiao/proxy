package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
)

func cmdGlobal(n int) {
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = SetGlobal(c); err != nil {
		fmt.Println(err)
		return
	}
	defer Reset()
	onExit(Reset)
	fmt.Println("proxy set to global mode")

	var cmd string
	if n < len(c.ProxyCommands) {
		cmd = c.ProxyCommands[n]
	}
	if cmd != "" {
		fmt.Println("start proxy")
		if err = StartProxy(cmd); err != nil {
			fmt.Println(err)
		}
	} else {
		<-(chan int)(nil)
	}
}

func onExit(f func() error) {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		f()
		os.Exit(0)
	}()
}

func cmdPAC(n int) {
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = SetPAC(c); err != nil {
		fmt.Println(err)
		return
	}
	defer Reset()
	onExit(Reset)
	fmt.Println("proxy set to pac mode")

	var addr = fmt.Sprintf("%s:%d", c.PACHost, c.PACPort)
	var cmd string
	if n < len(c.ProxyCommands) {
		cmd = c.ProxyCommands[n]
	}
	if cmd != "" {
		var ch = make(chan error, 2)

		fmt.Println("start proxy")
		go func() {
			ch <- StartProxy(cmd)
		}()

		fmt.Println("start pac server")
		go func() {
			ch <- StartServer(addr)
		}()
		// should not run to here, unless any of above goroutines exit.
		if err = <-ch; err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("start pac server")
		if err = StartServer(addr); err != nil {
			fmt.Println(err)
		}
	}
}

func cmdOff() {
	if err := Reset(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("proxy cleared")
	}
}

func cmdUpdate() {
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = UpdatePAC(c); err != nil {
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
