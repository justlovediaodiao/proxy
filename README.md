
# Proxy

Proxy is a tool to set os proxy settings, with support for [Proxy Auto Config](https://en.wikipedia.org/wiki/Proxy_auto-config).
It makes a pac file and set system proxy to pac/http/socks mode. **It is not a proxy server.**

Support Windows, MacOS and Ubuntu Desktop.

## Config

```
# resources/config.json
{
    "host": "127.0.0.1",       # proxy server host
    "port": 1080,              # proxy server port
    "protocol": "http",        # proxy protocol, http or socks/socks5
    "pac_host": "127.0.0.1",   # pac server listen host, default 127.0.0.1
    "pac_port": 1081,          # pac server listen port, default 1081
    "proxy_commands": [""],    # commands to start proxy process. you can add multi commands, optional
}
```

## Usage

```
proxy
    g[n]/global[n]: set os proxy to global mode. n: default 0, if 0 <= n < proxy_commands.length, start proxy_commands[n].
    pac[n]: set os proxy to pac mode.
    off/clear: clear os proxy setting.
    update: update pac file.
```

If you want to custom pac rule, write pac rule into `resources/user-rule.txt`.  
You need to run `update` command to update pac file when:
1. `host` or `port` config changed.
2. `resources/user-rule.txt` chaned.
3. you want to update [gfwlist](https://github.com/gfwlist/gfwlist).

## Build

- Command line version:
```bash
cd cmd
go build
```

- GUI version:
C Compiler is needed, see [https://developer.fyne.io/started/](https://developer.fyne.io/started/). Run following command based on your os to build:

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest

cd gui

# MacOS
fyne package -os darwin -icon icon.png

# Linux
fyne package -os linux -icon icon.png

# Windows
$env:CGO_ENABLED="1"    # Powershell
fyne package -os windows -icon icon.png
```

For windows, with no application icon, you can simply run following command to build:
```bash
cd gui
$env:CGO_ENABLED="1"
go build -ldflags -H=windowsgui
```
