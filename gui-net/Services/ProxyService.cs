using System.Diagnostics;
using System.Runtime.InteropServices;
using System.Text.Json;
using gui_net.Models;

namespace gui_net.Services;

public class ProxyService
{
    private Config? _config;
    private Process? _proxyProcess;
    private PacServer? _pacServer;

    public Config Config
    {
        get
        {
            if (_config == null)
                LoadConfig();
            return _config!;
        }
    }

    private void LoadConfig()
    {
        var configPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "resources", "config.json");
        if (File.Exists(configPath))
        {
            var json = File.ReadAllText(configPath);
            _config = JsonSerializer.Deserialize(json, JsonContext.Default.Config);
        }
        else
        {
            _config = new Config();
        }

        // Derive ProxyUrl
        if (_config != null)
        {
            _config.ProxyUrl = _config.Protocol switch
            {
                "http" => $"PROXY {_config.Host}:{_config.Port};DIRECT",
                "socks" => $"SOCKS://{_config.Host}:{_config.Port};DIRECT",
                "socks5" => $"SOCKS5://{_config.Host}:{_config.Port};DIRECT",
                _ => $"PROXY {_config.Host}:{_config.Port};DIRECT"
            };
        }
    }

    public void Off()
    {
        Reset();
        StopProxyProcess();
        StopPacServer();
    }

    public void Global()
    {
        Reset(); // Clear previous settings
        SetGlobal(Config);
        StartProxyProcess();
    }

    public void Pac()
    {
        Reset();
        SetPac(Config);
        StartPacServer();
        StartProxyProcess();
    }

    private void StartProxyProcess()
    {
        if (Config.ProxyCommands != null && Config.ProxyCommands.Length > 0 && _proxyProcess == null)
        {
            try
            {
                var cmd = Config.ProxyCommands[0];
                var parts = cmd.Split(' ', 2);
                var fileName = parts[0];
                var args = parts.Length > 1 ? parts[1] : "";

                var psi = new ProcessStartInfo
                {
                    FileName = fileName,
                    Arguments = args,
                    UseShellExecute = false,
                    CreateNoWindow = true
                };
                _proxyProcess = Process.Start(psi);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Failed to start proxy process: {ex.Message}");
            }
        }
    }

    private void StopProxyProcess()
    {
        if (_proxyProcess != null && !_proxyProcess.HasExited)
        {
            try
            {
                _proxyProcess.Kill();
            }
            catch { }
            _proxyProcess = null;
        }
    }

    private void StartPacServer()
    {
        _pacServer ??= new PacServer();
        _pacServer.Start(Config.PacHost, Config.PacPort);
    }

    private void StopPacServer()
    {
        if (_pacServer != null)
        {
            _pacServer.Stop();
            _pacServer = null;
        }
    }

    private void Reset()
    {
        if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
            WindowsProxy.Reset();
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.OSX))
            MacProxy.Reset();
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.Linux))
            LinuxProxy.Reset();
    }

    private void SetGlobal(Config config)
    {
        if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
            WindowsProxy.SetGlobal(config);
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.OSX))
            MacProxy.SetGlobal(config);
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.Linux))
            LinuxProxy.SetGlobal(config);
    }

    private void SetPac(Config config)
    {
        if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
            WindowsProxy.SetPac(config);
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.OSX))
            MacProxy.SetPac(config);
        else if (RuntimeInformation.IsOSPlatform(OSPlatform.Linux))
            LinuxProxy.SetPac(config);
    }
}