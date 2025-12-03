using System.Collections.ObjectModel;
using CommunityToolkit.Mvvm.ComponentModel;
using gui_net.Services;

namespace gui_net.ViewModels;

public partial class MainWindowViewModel : ObservableObject
{
    private readonly ProxyService _proxyService;
    [ObservableProperty]
    private string _selectedMode = "";

    [ObservableProperty]
    private string _statusMessage = "Ready";

    [ObservableProperty]
    private string _statusColor = "Black";

    public MainWindowViewModel()
    {
        _proxyService = new ProxyService();
        Modes = new ObservableCollection<string> { "Off", "Global", "Pac" };
    }

    public ObservableCollection<string> Modes { get; }

    partial void OnSelectedModeChanged(string value)
    {
        if (string.IsNullOrEmpty(value)) return;

        try
        {
            switch (value)
            {
                case "Off":
                    _proxyService.Off();
                    StatusMessage = "Proxy off";
                    StatusColor = "Green";
                    break;
                case "Global":
                    _proxyService.Global();
                    StatusMessage = "Global mode";
                    StatusColor = "Green";
                    break;
                case "Pac":
                    _proxyService.Pac();
                    StatusMessage = "Pac mode";
                    StatusColor = "Green";
                    break;
            }
        }
        catch (Exception ex)
        {
            StatusMessage = $"Error: {ex.Message}";
            StatusColor = "Red";
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
