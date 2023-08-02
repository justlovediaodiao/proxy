
from .config import get_config
from .pac import start_server
from .proxy import reset, set_global, set_pac, start_proxy


class Proxy:

    def __init__(self):
        self.config = get_config()
        self._proxy = None
        self._pac = None

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

    def global_(self):
        set_global(self.config)
        self._start_proxy()

    def pac(self):
        set_pac(self.config)
        self._start_pac()
        self._start_proxy()

    def off(self):
        reset()
        self._stop()