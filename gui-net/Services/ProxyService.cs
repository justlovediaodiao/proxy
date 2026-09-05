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

    public ProcessLogBuffer Logs { get; } = new();

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
        var configPath = "resources/config.json";
        if (!File.Exists(configPath))
            throw new FileNotFoundException($"Config file not found: {configPath}", configPath);


        var json = File.ReadAllText(configPath);
        _config = JsonSerializer.Deserialize(json, JsonContext.Default.Config);
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
        Reset();
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
        if (_proxyProcess is { HasExited: true })
        {
            _proxyProcess.Dispose();
            _proxyProcess = null;
        }

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
                    CreateNoWindow = true,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true
                };
                var process = new Process
                {
                    StartInfo = psi,
                    EnableRaisingEvents = true
                };
                process.OutputDataReceived += (_, e) =>
                {
                    if (e.Data != null)
                        Logs.Append("stdout", e.Data);
                };
                process.ErrorDataReceived += (_, e) =>
                {
                    if (e.Data != null)
                        Logs.Append("stderr", e.Data);
                };
                process.Exited += (_, _) =>
                {
                    try
                    {
                        Logs.Append("proxy", $"Process exited with code {process.ExitCode}.");
                    }
                    catch (InvalidOperationException)
                    {
                        Logs.Append("proxy", "Process exited.");
                    }
                };

                _proxyProcess = process;
                if (!process.Start())
                {
                    throw new InvalidOperationException("The proxy process could not be started.");
                }

                Logs.Append("proxy", $"Process started (PID {process.Id}).");
                process.BeginOutputReadLine();
                process.BeginErrorReadLine();
            }
            catch (Exception ex)
            {
                _proxyProcess?.Dispose();
                _proxyProcess = null;
                Logs.Append("proxy", $"Failed to start process: {ex.Message}");
            }
        }
    }

    private void StopProxyProcess()
    {
        if (_proxyProcess != null)
        {
            try
            {
                if (!_proxyProcess.HasExited)
                {
                    Logs.Append("proxy", "Stopping process.");
                    _proxyProcess.Kill();
                    _proxyProcess.WaitForExit(2_000);
                }
            }
            catch { }
            _proxyProcess.Dispose();
            _proxyProcess = null;
        }
    }

    private void StartPacServer()
    {
        _pacServer ??= new PacServer(Logs);
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
