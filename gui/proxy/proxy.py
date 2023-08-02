import subprocess
import sys

from .config import Config

if sys.platform == 'win32':

    def set_global(c: Config):
        if c.protocol == 'http':
            addr = f'{c.host}:{c.port}'
        else:
            addr = f'socks={c.host}:{c.port}'
        execute('resources/sysproxy.exe', 'global', addr)


    def set_pac(c: Config):
        addr = f'http://{c.host}:{c.port}/'
        execute('resources/sysproxy.exe', 'pac', addr)


    def reset():
        execute('resources/sysproxy.exe', 'set', '1', '-', '-', '-')


    def start_proxy(cmd: str) -> subprocess.Popen:
        return subprocess.Popen(' '.split(cmd), stdout=subprocess.PIPE, stderr=subprocess.PIPE)


    def execute(*args: str):
        subprocess.run(' '.join(args), stdout=subprocess.PIPE, stderr=subprocess.PIPE)
