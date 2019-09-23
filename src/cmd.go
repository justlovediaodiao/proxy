package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
)

func cmdGlobal() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if !p.Global {
		p.Global = true
		SaveConfig(p)
	}
	startGlobal(p)
}

func startPAC(p *Proxy) {
	var err = SetPAC(p)
	if err != nil {
		fmt.Println(err)
        return
	}
	var addr = fmt.Sprintf("%s:%d", p.PACHost, p.PACPort)
	fmt.Println("proxy set to pac mode\npac server is running...")
	// ctrl+c exit
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cmdOff()
		os.Exit(0)
	}()
	err = StartServer(addr)
	fmt.Println(err)
	cmdOff()
}

func startGlobal(p *Proxy) {
	var err = SetGlobal(p)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("proxy set to global mode")
	}
}

func cmdPAC() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if p.Global {
		p.Global = false
		SaveConfig(p)
	}
	startPAC(p)
}

func cmdOn() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if p.Global {
		startGlobal(p)
	} else {
		startPAC(p)
	}
}

func cmdOff() {
	var err = Reset()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("proxy cleared")
	}
}

func cmdUpdate() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	err = UpdatePAC(p.URL())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("pac updated")
	}
}

func cmdHelp() {
	fmt.Println("usage proxy [g/global/pac/on/off/update]")
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
	case "on":
		cmdOn()
	case "off":
		cmdOff()
	case "update":
		cmdUpdate()
	default:
		cmdHelp()
	}
}
