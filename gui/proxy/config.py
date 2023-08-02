import json
from dataclasses import dataclass


@dataclass
class Config:
    host: str
    port: int
    protocol: str
    pac_host: str
    pac_port: int
    proxy_commands: list[str]
    proxy_url: str


def get_config() -> Config:
    with open('resources/config.json') as fp:
        content = fp.read()
    c = json.loads(content)
    match c['protocol']:
        case 'http':
            proxy_url = 'PROXY %s%s;DIRECT'
        case 'socks':
            proxy_url = 'SOCKS://%s%s;DIRECT'
        case 'socks5':
            proxy_url = 'SOCKS5://%s%s;DIRECT'
        case _:
            raise ValueError('Unknown proxy protocol')
    c['proxy_url'] = proxy_url % (c['host'], c['port'])
    return Config(**c)
