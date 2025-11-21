import ctypes
from enum import IntEnum


class InternetOption(IntEnum):
    Refresh = 37
    PerConnectionOption = 75
    SettingsChanged = 39
    ProxySettingsChanged = 95


class InternetPerConnOption(IntEnum):
    Flags = 1
    ProxyServer = 2
    ProxyBypass = 3
    AutoConfigUrl = 4


class InternetPerConnFlags(IntEnum):
    Direct = 1
    Proxy = 2
    AutoProxyUrl = 4
    AutoDetect = 8


class InternetPerConnOptionList(ctypes.Structure):
    _fields_ = [
        ("dwSize", ctypes.c_uint32),
        ("pszConnection", ctypes.c_char_p),
        ("dwOptionCount", ctypes.c_uint32),
        ("dwOptionError", ctypes.c_uint32),
        ("pOptions", ctypes.c_void_p)
    ]


class InternetConnOption(ctypes.Structure):
    _fields_ = [
        ("dwOption", ctypes.c_uint32),
        ("dwValue", ctypes.c_void_p)  # DWORD | LPSTR
    ]


def reset():
    op = InternetConnOption(
        dwOption=InternetPerConnOption.Flags,
        dwValue=InternetPerConnFlags.Direct
    )

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=1,
        dwOptionError=0,
        pOptions=ctypes.cast(ctypes.pointer(op), ctypes.c_void_p)
    )

    _set_system_proxy(opl)


def set_global(proxy: str, bypass: str):
    op0 = InternetConnOption(
        dwOption=InternetPerConnOption.Flags,
        dwValue=InternetPerConnFlags.Proxy | InternetPerConnFlags.Direct
    )

    op1 = InternetConnOption(
        dwOption=InternetPerConnOption.ProxyServer,
        dwValue=ctypes.cast(ctypes.c_char_p(proxy.encode()), ctypes.c_void_p)
    )

    op2 = InternetConnOption(
        dwOption=InternetPerConnOption.ProxyBypass,
        dwValue=ctypes.cast(ctypes.c_char_p(bypass.encode()), ctypes.c_void_p)
    )

    ops = (InternetConnOption * 3)(op0, op1, op2)

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=3,
        dwOptionError=0,
        pOptions=ctypes.cast(ctypes.pointer(ops), ctypes.c_void_p)
    )

    _set_system_proxy(opl)


def set_pac(proxy_url: str):
    op0 = InternetConnOption(
        dwOption=InternetPerConnOption.Flags,
        dwValue=InternetPerConnFlags.AutoProxyUrl | InternetPerConnFlags.Direct
    )

    op1 = InternetConnOption(
        dwOption=InternetPerConnOption.AutoConfigUrl,
        dwValue=ctypes.cast(ctypes.c_char_p(proxy_url.encode()), ctypes.c_void_p)
    )

    ops = (InternetConnOption * 2)(op0, op1)

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=2,
        dwOptionError=0,
        pOptions=ctypes.cast(ctypes.pointer(ops), ctypes.c_void_p)
    )

    _set_system_proxy(opl)


def _set_system_proxy(opl: InternetPerConnOptionList):
    internet_set_option = ctypes.windll.Wininet.InternetSetOptionA
    r = internet_set_option(None, InternetOption.PerConnectionOption, ctypes.pointer(opl), ctypes.sizeof(opl))
    _handle_error(r)
    r = internet_set_option(None, InternetOption.ProxySettingsChanged, None, 0)
    _handle_error(r)
    r = internet_set_option(None, InternetOption.Refresh, None, 0)
    _handle_error(r)


def _handle_error(r: int):
    if r == 1:
        return
    code = ctypes.get_last_error()
    err = ctypes.FormatError(code)
    e = OSError()
    e.errno = code
    e.winerror = code
    e.strerror = err
    raise e
