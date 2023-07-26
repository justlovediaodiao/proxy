package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/justlovediaodiao/proxy"
)

var (
	pacStarted    = false
	proxyProceess *os.Process
	tipContent    *widget.Label
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
	if proxyProceess != nil {
		proxyProceess.Signal(os.Interrupt)
		proxyProceess = nil
	}
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
	if proxyProceess == nil && len(c.ProxyCommands) != 0 {
		p, err := proxy.StartProxyAsync(c.ProxyCommands[0])
		if err != nil {
			show("start proxy server error", err)
			return
		}
		proxyProceess = p
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

	if proxyProceess == nil && len(c.ProxyCommands) != 0 {
		p, err := proxy.StartProxyAsync(c.ProxyCommands[0])
		if err != nil {
			show("start proxy server error", err)
			return
		}
		proxyProceess = p
	}
	show("pac mode", nil)
}
