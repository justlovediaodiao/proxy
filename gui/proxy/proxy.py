import subprocess
import sys

from .config import Config

if sys.platform == 'win32':

    def set_global(c: Config):
        if c.protocol == 'http':
            addr = f'{c.host}:{c.port}'
        else:
            addr = f'socks={c.host}:{c.port}'
        execute('resources/sysproxy.exe', 'global', addr, '<local>;192.168.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*')

    def set_pac(c: Config):
        addr = f'http://{c.host}:{c.port}/'
        execute('resources/sysproxy.exe', 'pac', addr)

    def reset():
        execute('resources/sysproxy.exe', 'set', '1', '-', '-', '-')

    def start_proxy(cmd: str) -> subprocess.Popen:
        return subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, creationflags=subprocess.CREATE_NO_WINDOW)

    def execute(*args: str):
        subprocess.run(' '.join(args), stdout=subprocess.PIPE, stderr=subprocess.PIPE, creationflags=subprocess.CREATE_NO_WINDOW)