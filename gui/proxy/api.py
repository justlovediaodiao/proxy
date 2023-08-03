
from .config import get_config
from .pac import start_server
from .proxy import reset, set_global, set_pac, start_proxy


class Proxy:

    def __init__(self):
        self._config = None
        self._proxy = None
        self._pac = None
        self._state = ''

    @property
    def config(self):
        if not self._config:
            self._config = get_config()
        return self._config

    def _start_proxy(self):
        if self.config.proxy_commands and not self._proxy:
            cmd = self.config.proxy_commands[0]
            self._proxy = start_proxy(cmd)

    def _start_pac(self):
        if not self._pac:
            self._pac = start_server(self.config.pac_host, self.config.pac_port)

    def _stop(self):
        if self._proxy:
            self._proxy.kill()
            self._proxy = None
        if self._pac:
            self._pac.shutdown()
            self._pac = None

    def _reset(self):
        if self._state != 'off':
            reset()
            self._state = 'off'

    def global_(self):
        self._reset()
        set_global(self.config)
        self._start_proxy()

    def pac(self):
        self._reset()
        set_pac(self.config)
        self._start_pac()
        self._start_proxy()

    def off(self):
        self._reset()