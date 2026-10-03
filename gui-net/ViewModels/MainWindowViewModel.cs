using CommunityToolkit.Mvvm.ComponentModel;
using gui_net.Services;
using Avalonia.Threading;

namespace gui_net.ViewModels;

public partial class MainWindowViewModel : ObservableObject
{
    private readonly ProxyService _proxyService;

    private bool _isApplying;
    private int _applyVersion;

    [ObservableProperty]
    private bool _proxyEnabled;

    [ObservableProperty]
    [NotifyPropertyChangedFor(nameof(PacModeSelected))]
    private bool _globalModeSelected = true;

    public bool PacModeSelected => !GlobalModeSelected;

    [ObservableProperty]
    private bool _controlsEnabled = true;

    [ObservableProperty]
    private bool _isApplyingStatus;

    [ObservableProperty]
    private string _statusMessage = "Off";

    [ObservableProperty]
    private string _statusColor = "#8A8A8A";

    [ObservableProperty]
    private string? _statusDetail;

    public MainWindowViewModel()
    {
        _proxyService = new ProxyService();
        _proxyService.ProcessExitedUnexpectedly += OnProxyProcessExited;
    }

    private void OnProxyProcessExited()
    {
        Dispatcher.UIThread.Post(() =>
        {
            if (ProxyEnabled && _proxyService.ProcessFailure is { } message)
                ShowError(message);
        });
    }

    public bool HasError => StatusMessage == "Error";

    partial void OnStatusMessageChanged(string value) => OnPropertyChanged(nameof(HasError));

    public bool HasStatusDetail => !string.IsNullOrWhiteSpace(StatusDetail);

    public ProcessLogBuffer Logs => _proxyService.Logs;

    partial void OnStatusDetailChanged(string? value)
    {
        OnPropertyChanged(nameof(HasStatusDetail));
    }

    partial void OnProxyEnabledChanged(bool value)
    {
        RequestApply();
    }

    partial void OnGlobalModeSelectedChanged(bool value)
    {
        if (ProxyEnabled)
            RequestApply();
    }

    private void RequestApply()
    {
        _applyVersion++;

        if (_isApplying)
            return;

        _ = ApplyRequestedModeAsync();
    }

    private async Task ApplyRequestedModeAsync()
    {
        _isApplying = true;
        IsApplyingStatus = true;
        ControlsEnabled = false;

        try
        {
            int handledVersion;
            do
            {
                handledVersion = _applyVersion;
                await ApplyCurrentModeAsync();
            }
            while (handledVersion != _applyVersion);
        }
        finally
        {
            ControlsEnabled = true;
            IsApplyingStatus = false;
            _isApplying = false;
        }
    }

    private async Task ApplyCurrentModeAsync()
    {
        var enabled = ProxyEnabled;
        var mode = PacModeSelected ? "Pac" : "Global";

        StatusDetail = null;
        StatusColor = "#D9A300";
        StatusMessage = String.Empty;

        try
        {
            await Task.Run(() =>
            {
                if (!enabled)
                {
                    _proxyService.Off();
                    return;
                }

                if (mode == "Pac")
                    _proxyService.Pac();
                else
                    _proxyService.Global();
            });

            if (enabled && _proxyService.ProcessFailure is { } failure)
            {
                ShowError(failure);
                return;
            }

            StatusColor = enabled ? "#16A34A" : "#8A8A8A";
            StatusMessage = enabled ? "On" : "Off";
        }
        catch (Exception e)
        {
            ShowError(e.Message);
        }
    }

    private void ShowError(string message)
    {
        StatusMessage = "Error";
        StatusColor = "#D13438";
        StatusDetail = message;
    }

    public void OnExit()
    {
        _proxyService.ProcessExitedUnexpectedly -= OnProxyProcessExited;
        try
        {
            _proxyService.Off();
        }
        catch { }
    }
}
