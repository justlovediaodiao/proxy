
### proxy

proxy is a tool to set os proxy setting, with support for [Proxy Auto Config](https://en.wikipedia.org/wiki/Proxy_auto-config).
it makes a pac file and set system proxy to pac/http/socks mode. **it is not a proxy server.**


### config

```
# resources/config.json
{
    "host": "127.0.0.1",    # proxy server host
    "port": 1080,           # proxy server port
    "protocol": "http",     # proxy protocol, http or socks/socks5
    "pac_host": "127.0.0.1",# pac server listen host, default 127.0.0.1
    "pac_port": 1081,       # pac server listen port, default 1081
    "proxy_command": "",    # command to start proxy process, optional
}
```

### usage

```
proxy
    g/global: set os proxy to global mode.
    pac: set os proxy to pac mode.
    off/clear: clear os proxy setting.
    update: update pac file.
```

if you want to custom pac rule, write pac rule into `resources/user-rule.txt`.  
you need to run `update` command to update pac file when:
1. `host` or `port` config changed.
2. `resources/user-rule.txt` chaned.
3. you want to update [gfwlist](https://github.com/gfwlist/gfwlist).

