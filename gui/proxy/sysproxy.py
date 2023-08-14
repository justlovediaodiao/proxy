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
        dwValue=InternetPerConnFlags.AutoDetect | InternetPerConnFlags.Direct
    )

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=1,
        dwOptionError=0,
        pOptions=ctypes.pointer(op)
    )

    _set_system_proxy(opl)


def set_global(proxy: str, bypass: str):
    op0 = InternetConnOption(
        dwOption=InternetPerConnOption.Flags,
        dwValue=InternetPerConnFlags.Proxy | InternetPerConnFlags.Direct
    )

    op1 = InternetConnOption(
        dwOption=InternetPerConnOption.ProxyServer,
        dwValue=ctypes.c_char_p(proxy.encode())
    )

    op2 = InternetConnOption(
        dwOption=InternetPerConnOption.ProxyBypass,
        dwValue=ctypes.c_char_p(bypass.encode())
    )

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=3,
        dwOptionError=0,
        pOptions=ctypes.pointer((InternetConnOption * 3)(op0, op1, op2))
    )

    _set_system_proxy(opl)


def set_pac(proxy_url: str):
    op0 = InternetConnOption(
        dwOption=InternetPerConnOption.Flags,
        dwValue=InternetPerConnFlags.AutoProxyUrl | InternetPerConnFlags.Direct
    )

    op1 = InternetConnOption(
        dwOption=InternetPerConnOption.AutoConfigUrl,
        dwValue=ctypes.c_char_p(proxy_url.encode())
    )

    opl = InternetPerConnOptionList(
        dwSize=ctypes.sizeof(InternetPerConnOptionList),
        pszConnection=None,
        dwOptionCount=2,
        dwOptionError=0,
        pOptions=ctypes.pointer((InternetConnOption * 2)(op0, op1))
    )

    _set_system_proxy(opl)


def _set_system_proxy(opl: InternetPerConnOptionList):
    internet_set_option = ctypes.windll.Wininet.InternetSetOptionA
    internet_set_option(None, InternetOption.PerConnectionOption, opl, ctypes.sizeof(opl))
    internet_set_option(None, InternetOption.ProxySettingsChanged, None, 0)
    internet_set_option(None, InternetOption.Refresh, None, 0)
