using System.Diagnostics;
using System.Runtime.InteropServices;
using gui_net.Models;

namespace gui_net.Services;

internal static class MacProxy
{
    public static void SetGlobal(Config c)
    {
        var networks = ListNetworks();
        foreach (var network in networks)
        {
            if (c.Protocol == "http")
            {
                Execute("networksetup", $"-setwebproxy \"{network}\" {c.Host} {c.Port}");
                Execute("networksetup", $"-setsecurewebproxy \"{network}\" {c.Host} {c.Port}");
            }
            else if (c.Protocol.StartsWith("socks"))
            {
                Execute("networksetup", $"-setsocksfirewallproxy \"{network}\" {c.Host} {c.Port}");
            }
        }
    }

    public static void SetPac(Config c)
    {
        var networks = ListNetworks();
        var url = $"http://{c.PacHost}:{c.PacPort}/";
        foreach (var network in networks)
        {
            Execute("networksetup", $"-setautoproxyurl \"{network}\" \"{url}\"");
        }
    }

    public static void Reset()
    {
        var networks = ListNetworks();
        foreach (var network in networks)
        {
            Execute("networksetup", $"-setautoproxystate \"{network}\" off");
            Execute("networksetup", $"-setwebproxystate \"{network}\" off");
            Execute("networksetup", $"-setsecurewebproxystate \"{network}\" off");
            Execute("networksetup", $"-setsocksfirewallproxystate \"{network}\" off");
        }
    }

    private static List<string> ListNetworks()
    {
        var output = ExecuteAndGetOutput("networksetup", "-listallnetworkservices");
        var networks = new List<string>();
        foreach (var line in output.Split('\n'))
        {
            if (line.Contains("Wi-Fi") || line.Contains("Ethernet") || line.Contains("Thunderbolt"))
            {
                if (line.Contains("*")) continue; // Skip disabled
                if (line.Trim() == "Wi-Fi" || line.Trim() == "Ethernet")
                {
                    networks.Add(line.Trim());
                }
            }
        }
        return networks;
    }

    private static void Execute(string cmd, string args)
    {
        Process.Start(new ProcessStartInfo
        {
            FileName = cmd,
            Arguments = args,
            UseShellExecute = false,
            CreateNoWindow = true
        })?.WaitForExit();
    }

    private static string ExecuteAndGetOutput(string cmd, string args)
    {
        var psi = new ProcessStartInfo
        {
            FileName = cmd,
            Arguments = args,
            RedirectStandardOutput = true,
            UseShellExecute = false,
            CreateNoWindow = true
        };
        var p = Process.Start(psi);
        if (p == null) return string.Empty;
        var output = p.StandardOutput.ReadToEnd();
        p.WaitForExit();
        return output;
    }
}

internal static class LinuxProxy
{
    public static void SetGlobal(Config c)
    {
        if (c.Protocol == "http")
        {
            Execute("gsettings", $"set org.gnome.system.proxy.http host '{c.Host}'");
            Execute("gsettings", $"set org.gnome.system.proxy.http port {c.Port}");
            Execute("gsettings", $"set org.gnome.system.proxy.https host '{c.Host}'");
            Execute("gsettings", $"set org.gnome.system.proxy.https port {c.Port}");
            Execute("gsettings", "set org.gnome.system.proxy mode 'manual'");
        }
        else if (c.Protocol.StartsWith("socks"))
        {
            Execute("gsettings", $"set org.gnome.system.proxy.socks host '{c.Host}'");
            Execute("gsettings", $"set org.gnome.system.proxy.socks port {c.Port}");
            Execute("gsettings", "set org.gnome.system.proxy mode 'manual'");
        }
    }

    public static void SetPac(Config c)
    {
        var url = $"http://{c.PacHost}:{c.PacPort}/";
        Execute("gsettings", $"set org.gnome.system.proxy autoconfig-url '{url}'");
        Execute("gsettings", "set org.gnome.system.proxy mode 'auto'");
    }

    public static void Reset()
    {
        Execute("gsettings", "set org.gnome.system.proxy mode 'none'");
    }

    private static void Execute(string cmd, string args)
    {
        Process.Start(new ProcessStartInfo
        {
            FileName = cmd,
            Arguments = args,
            UseShellExecute = false,
            CreateNoWindow = true
        })?.WaitForExit();
    }
}

internal static class WindowsProxy
{
    // P/Invoke definitions
    [DllImport("wininet.dll", SetLastError = true, CharSet = CharSet.Ansi)]
    private static extern bool InternetSetOption(IntPtr hInternet, int dwOption, IntPtr lpBuffer, int dwBufferLength);

    private const int INTERNET_OPTION_PER_CONNECTION_OPTION = 75;
    private const int INTERNET_OPTION_SETTINGS_CHANGED = 39;
    private const int INTERNET_OPTION_REFRESH = 37;

    private const int PROXY_TYPE_DIRECT = 0x00000001;
    private const int PROXY_TYPE_PROXY = 0x00000002;
    private const int PROXY_TYPE_AUTO_PROXY_URL = 0x00000004;

    private const int INTERNET_PER_CONN_FLAGS = 1;
    private const int INTERNET_PER_CONN_PROXY_SERVER = 2;
    private const int INTERNET_PER_CONN_PROXY_BYPASS = 3;
    private const int INTERNET_PER_CONN_AUTOCONFIG_URL = 4;

    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Ansi)]
    private struct INTERNET_PER_CONN_OPTION_LIST
    {
        public int dwSize;
        public IntPtr pszConnection;
        public int dwOptionCount;
        public int dwOptionError;
        public IntPtr pOptions;
    }

    [StructLayout(LayoutKind.Explicit)]
    private struct INTERNET_PER_CONN_OPTION_OPTION_UNION
    {
        [FieldOffset(0)]
        public int dwValue;
        [FieldOffset(0)]
        public IntPtr pszValue;
        [FieldOffset(0)]
        public System.Runtime.InteropServices.ComTypes.FILETIME ftValue;
    }

    [StructLayout(LayoutKind.Sequential)]
    private struct INTERNET_PER_CONN_OPTION
    {
        public int dwOption;
        public INTERNET_PER_CONN_OPTION_OPTION_UNION Value;
    }

    public static void SetGlobal(Config c)
    {
        var addr = c.Protocol == "http"
            ? $"{c.Host}:{c.Port}"
            : $"socks={c.Host}:{c.Port}";

        var options = new INTERNET_PER_CONN_OPTION[3];

        // Flags
        options[0].dwOption = INTERNET_PER_CONN_FLAGS;
        options[0].Value.dwValue = PROXY_TYPE_PROXY | PROXY_TYPE_DIRECT;

        // Proxy Server
        options[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER;
        options[1].Value.pszValue = Marshal.StringToHGlobalAnsi(addr);

        // Bypass
        options[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS;
        options[2].Value.pszValue = Marshal.StringToHGlobalAnsi("<local>;192.168.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*");

        ApplyOptions(options);

        Marshal.FreeHGlobal(options[1].Value.pszValue);
        Marshal.FreeHGlobal(options[2].Value.pszValue);
    }

    public static void SetPac(Config c)
    {
        var url = $"http://{c.PacHost}:{c.PacPort}/";
        var options = new INTERNET_PER_CONN_OPTION[2];

        // Flags
        options[0].dwOption = INTERNET_PER_CONN_FLAGS;
        options[0].Value.dwValue = PROXY_TYPE_AUTO_PROXY_URL | PROXY_TYPE_DIRECT;

        // Auto Config URL
        options[1].dwOption = INTERNET_PER_CONN_AUTOCONFIG_URL;
        options[1].Value.pszValue = Marshal.StringToHGlobalAnsi(url);

        ApplyOptions(options);

        Marshal.FreeHGlobal(options[1].Value.pszValue);
    }

    public static void Reset()
    {
        var options = new INTERNET_PER_CONN_OPTION[1];
        options[0].dwOption = INTERNET_PER_CONN_FLAGS;
        options[0].Value.dwValue = PROXY_TYPE_DIRECT;

        ApplyOptions(options);
    }

    private static void ApplyOptions(INTERNET_PER_CONN_OPTION[] options)
    {
        var optionSize = Marshal.SizeOf<INTERNET_PER_CONN_OPTION>();
        var optionsPtr = Marshal.AllocCoTaskMem(optionSize * options.Length);

        for (var i = 0; i < options.Length; i++)
        {
            var ptr = new IntPtr(optionsPtr.ToInt64() + (i * optionSize));
            Marshal.StructureToPtr(options[i], ptr, false);
        }

        var list = new INTERNET_PER_CONN_OPTION_LIST
        {
            dwSize = Marshal.SizeOf<INTERNET_PER_CONN_OPTION_LIST>(),
            pszConnection = IntPtr.Zero,
            dwOptionCount = options.Length,
            dwOptionError = 0,
            pOptions = optionsPtr
        };

        var listSize = Marshal.SizeOf<INTERNET_PER_CONN_OPTION_LIST>();
        var listPtr = Marshal.AllocCoTaskMem(listSize);
        Marshal.StructureToPtr(list, listPtr, false);

        InternetSetOption(IntPtr.Zero, INTERNET_OPTION_PER_CONNECTION_OPTION, listPtr, listSize);
        
        // Notify system
        InternetSetOption(IntPtr.Zero, INTERNET_OPTION_SETTINGS_CHANGED, IntPtr.Zero, 0);
        InternetSetOption(IntPtr.Zero, INTERNET_OPTION_REFRESH, IntPtr.Zero, 0);

        Marshal.FreeCoTaskMem(optionsPtr);
        Marshal.FreeCoTaskMem(listPtr);
    }
}
