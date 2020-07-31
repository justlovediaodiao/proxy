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
	startGlobal(p)
}

func startPAC(p *Proxy) {
	var err = SetPAC(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	var addr = fmt.Sprintf("%s:%d", p.PACHost, p.PACPort)
	fmt.Println("proxy set to pac mode")
	onClose(cmdOff)
	if p.ProxyCommand != "" {
		go startProxy(p)
	}
	fmt.Println("start pac server")
	err = StartServer(addr)
	fmt.Println(err)
	cmdOff()
}

func startGlobal(p *Proxy) {
	var err = SetGlobal(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("proxy set to global mode")
	onClose(cmdOff)
	if p.ProxyCommand != "" {
		go startProxy(p)
	}
	<-(chan int)(nil)
}

func onClose(handler func()) {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		handler()
		os.Exit(0)
	}()
}

func startProxy(p *Proxy) {
	fmt.Println("start proxy")
	err := StartProxy(p.ProxyCommand)
	if err != nil {
		fmt.Println(err)
	}
}

func cmdPAC() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	startPAC(p)
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
	err = UpdatePAC(p)
	if err != nil {
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
