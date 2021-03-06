package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
)

func cmdGlobal() {
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if err = SetGlobal(c); err != nil {
		fmt.Println(err)
		return
	}
	defer Reset()
	onExit(Reset)
	fmt.Println("proxy set to global mode")

	if c.ProxyCommand != "" {
		fmt.Println("start proxy")
		if err = StartProxy(c.ProxyCommand); err != nil {
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

func cmdPAC() {
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if err = SetPAC(c); err != nil {
		fmt.Println(err)
		return
	}
	defer Reset()
	onExit(Reset)
	fmt.Println("proxy set to pac mode")

	var addr = fmt.Sprintf("%s:%d", c.PACHost, c.PACPort)
	if c.ProxyCommand != "" {
		var ch = make(chan error, 2)

		fmt.Println("start proxy")
		go func() {
			ch <- StartProxy(c.ProxyCommand)
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
	}
	if err = UpdatePAC(c); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("pac updated")
	}
}

func cmdHelp() {
	fmt.Println("usage proxy [g/global/pac/off/clear/update]")
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
	switch cmd {
	case "g", "global":
		cmdGlobal()
	case "pac":
		cmdPAC()
	case "off", "clear":
		cmdOff()
	case "update":
		cmdUpdate()
	default:
		cmdHelp()
	}
}
