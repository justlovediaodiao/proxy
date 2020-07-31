
### proxy

some proxies do not support pac([Proxy Auto Config](https://en.wikipedia.org/wiki/Proxy_auto-config)) mode, so this tool comes.    
it makes a pac file and set system proxy to pac mode. **it is not a proxy server.**

### build

run `go build src`

### config

```
# resources/config.json
{
    "host": "127.0.0.1",    # proxy server host
    "port": 1080,           # proxy server port
    "protocol": "http",     # proxy protocol, http or socks/socks5/socks4
    "pac_host": "127.0.0.1",# pac server listen host, default 127.0.0.1
    "pac_port": 1081,       # pac server listen port, default 1081
    "proxy_command": "",    # command to start proxy process
}
```

### usage

supported commands list:
```
proxy
    g/global: set proxy to global mode.
    pac: set proxy to pac mode.
    off/clear: clear proxy.
    update: update pac file.
```

if you want to custom pac rule, write pac rule into `resources/user-rule.txt`.  
macos and windows10 do not support local pac file, so the tool use a built-in http server to return pac file. when you set proxy to pac mode, the server will start. proxy will be reset when pac server exit.

**important** 
- the executable file must be in the same directory with `resources` directory.  
- you need to run `update` command to update pac file when:
1. `host` or `port` config changed.
2. `resources/user-rule.txt` chaned.
3. you want to update [gfwlist](https://github.com/gfwlist/gfwlist).

  