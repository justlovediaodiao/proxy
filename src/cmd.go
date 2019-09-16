package main

import (
	"fmt"
	"os"
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
	err = SetGlobal(p)
	if err != nil {
		fmt.Println(err)
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
	err = SetPAC(p)
	if err != nil {
		fmt.Println(err)
	}
	var addr = fmt.Sprintf("%s:%d", p.PACHost, p.PACPort)
	StartServer(addr)
	cmdOff()
}

func cmdOn() {
	p, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	if p.Global {
		err = SetGlobal(p)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = SetPAC(p)
		if err != nil {
			fmt.Println(err)
		}
		var addr = fmt.Sprintf("%s:%d", p.PACHost, p.PACPort)
		StartServer(addr)
		cmdOff()
	}
}

func cmdOff() {
	var err = Reset()
	if err != nil {
		fmt.Println(err)
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
