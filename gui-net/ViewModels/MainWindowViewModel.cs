using CommunityToolkit.Mvvm.ComponentModel;
using Avalonia.Threading;
using gui_net.Services;

namespace gui_net.ViewModels;

public partial class MainWindowViewModel : ObservableObject
{
    private readonly ProxyService _proxyService;
    private readonly DispatcherTimer _spinnerTimer;

    private bool _isApplying;
    private int _applyVersion;

    [ObservableProperty]
    private bool _proxyEnabled;

    [ObservableProperty]
    private bool _pacModeSelected = true;

    [ObservableProperty]
    private bool _globalModeSelected;

    [ObservableProperty]
    private bool _controlsEnabled = true;

    [ObservableProperty]
    private bool _isApplyingStatus;

    [ObservableProperty]
    private double _spinnerAngle;

    [ObservableProperty]
    private string _statusMessage = "Off";

    [ObservableProperty]
    private string _statusColor = "#8A8A8A";

    [ObservableProperty]
    private string? _statusDetail;

    public MainWindowViewModel()
    {
        _proxyService = new ProxyService();
        _spinnerTimer = new DispatcherTimer
        {
            Interval = TimeSpan.FromMilliseconds(16)
        };
        _spinnerTimer.Tick += (_, _) => SpinnerAngle = (SpinnerAngle + 8) % 360;
    }

    public bool HasStatusDetail => !string.IsNullOrWhiteSpace(StatusDetail);

    partial void OnStatusDetailChanged(string? value)
    {
        OnPropertyChanged(nameof(HasStatusDetail));
    }

    partial void OnProxyEnabledChanged(bool value)
    {
        RequestApply();
    }

    partial void OnPacModeSelectedChanged(bool value)
    {
        if (value && ProxyEnabled)
            RequestApply();
    }

    partial void OnGlobalModeSelectedChanged(bool value)
    {
        if (value && ProxyEnabled)
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
        _spinnerTimer.Start();
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
            _spinnerTimer.Stop();
            SpinnerAngle = 0;
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

            StatusColor = enabled ? "#16A34A" : "#8A8A8A";
            StatusMessage = enabled ? "On" : "Off";
        }
        catch (Exception e)
        {
            StatusMessage = "Error";
            StatusColor = "#D13438";
            StatusDetail = e.Message;
        }
    }

    public void OnExit()
    {
        try
        {
            _proxyService.Off();
        }
        catch { }
    }
}
