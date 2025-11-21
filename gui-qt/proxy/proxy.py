import subprocess
import sys

from .config import Config


_creation_flags = 0

if sys.platform == 'win32':

    from . import sysproxy

    _creation_flags = subprocess.CREATE_NO_WINDOW

    def set_global(c: Config):
        if c.protocol == 'http':
            addr = f'{c.host}:{c.port}'
        elif c.protocol.startswith('socks'):
            addr = f'socks={c.host}:{c.port}'
        
        sysproxy.set_global(addr, '<local>;192.168.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*')

    def set_pac(c: Config):
        addr = f'http://{c.host}:{c.port}/'
        sysproxy.set_pac(addr)

    def reset():
        sysproxy.reset()


if sys.platform == 'darwin':

    def set_global(c: Config):
        networks = list_network()
        for network in networks:
            if c.protocol == 'http':
                execute('networksetup', '-setwebproxy', network, c.host, c.port)
                execute('networksetup', '-setsecurewebproxy', network, c.host, c.port)
            elif c.protocol.startswith('socks'):
                execute('networksetup', '-setsocksfirewallproxy', network, c.host, c.port)

    def set_pac(c: Config):
        networks = list_network()
        url = f'http://{c.host}:{c.port}/'
        for network in networks:
            execute('networksetup', '-setautoproxyurl', network, url)

    def reset():
        networks = list_network()
        for network in networks:
            execute('networksetup', '-setautoproxystate', network, 'off')
            execute('networksetup', '-setwebproxystate', network, 'off')
            execute('networksetup', '-setsecurewebproxystate', network, 'off')
            execute('networksetup', '-setsocksfirewallproxystate', network, 'off')

    def list_network() -> list[str]:
        result = subprocess.run('networksetup', '-listallnetworkservices').stdout
        networks = []
        for network in result.decode().split('\n'):
            if network == 'Wi-Fi' or network == 'Ethernet':
                networks.append(network)
        return networks


if sys.platform == 'linux':

    def set_global(c: Config):
        if c.protocol == 'http':
            execute('gsettings', 'set', 'org.gnome.system.proxy.http', 'host', c.host)
            execute('gsettings', 'set', 'org.gnome.system.proxy.http', 'port', c.port)
            execute('gsettings', 'set', 'org.gnome.system.proxy.https', 'host', c.host)
            execute('gsettings', 'set', 'org.gnome.system.proxy.https', 'port', c.port)
            execute('gsettings', 'set', 'org.gnome.system.proxy', 'mode', 'manual')
        elif c.protocol.startswith('socks'):
            execute('gsettings', 'set', 'org.gnome.system.proxy.socks', 'host', c.host)
            execute('gsettings', 'set', 'org.gnome.system.proxy.socks', 'port', c.port)
            execute('gsettings', 'set', 'org.gnome.system.proxy', 'mode', 'manual')

    def set_pac(c: Config):
        url = f'http://{c.host}:{c.port}/'
        execute('gsettings', 'set', 'org.gnome.system.proxy', 'autoconfig-url', url)
        execute('gsettings', 'set', 'org.gnome.system.proxy', 'mode', 'auto')

    def reset():
        execute('gsettings', 'set', 'org.gnome.system.proxy', 'mode', 'none')


def execute(*args: str):
        subprocess.run(' '.join(args), creationflags=_creation_flags)

def start_proxy(cmd: str) -> subprocess.Popen:
    return subprocess.Popen(cmd, creationflags=_creation_flags)
