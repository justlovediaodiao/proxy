package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/justlovediaodiao/proxy"
)

var (
	pacStarted   = false
	proxyStarted = false
	tipContent   *widget.Label
)

func main() {
	w := app.New().NewWindow("Proxy")

	combo := widget.NewSelect([]string{"Off", "Global", "Pac"}, func(value string) {
		switch value {
		case "Off":
			off()
		case "Global":
			global()
		case "Pac":
			pac()
		}
	})

	tipContent = widget.NewLabel("")

	w.SetContent(container.NewVBox(combo, tipContent))
	w.Resize(fyne.Size{Width: 200, Height: 160})

	w.SetCloseIntercept(func() {
		off()
		w.Close()
	})
	w.ShowAndRun()
}

func show(tip string, e error) {
	if e == nil {
		tipContent.SetText(tip)
	} else {
		tipContent.SetText(tip + ": " + e.Error())
	}
}

func off() {
	if err := proxy.Reset(); err != nil {
		show("reset system proxy error", err)
		return
	}
	show("proxy off", nil)
}

func global() {
	c, err := proxy.GetConfig()
	if err != nil {
		show("get config error", err)
		return
	}
	if err = proxy.SetGlobal(c); err != nil {
		show("set system proxy error", err)
		return
	}
	if !proxyStarted && len(c.ProxyCommands) != 0 {
		proxyStarted = true
		go func() {
			if err := proxy.StartProxy(c.ProxyCommands[0]); err != nil {
				log.Printf("start proxy server error: %s", err)
			}
		}()
	}
	show("global mode", nil)
}

func pac() {
	c, err := proxy.GetConfig()
	if err != nil {
		show("get config error", err)
		return
	}
	if err = proxy.SetPAC(c); err != nil {
		show("set system proxy error", err)
		return
	}
	var addr = fmt.Sprintf("%s:%d", c.PACHost, c.PACPort)

	if !pacStarted {
		pacStarted = true
		go func() {
			if err := proxy.StartServer(addr); err != nil {
				log.Printf("start pac server error: %s", err)
			}
		}()
	}

	if !proxyStarted && len(c.ProxyCommands) != 0 {
		proxyStarted = true
		go func() {
			if err := proxy.StartProxy(c.ProxyCommands[0]); err != nil {
				log.Printf("start proxy server error: %s", err)
			}
		}()
	}
	show("pac mode", nil)
}
