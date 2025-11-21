using System;
using System.Collections.ObjectModel;
using ReactiveUI;
using gui_net.Services;

namespace gui_net.ViewModels
{
    public class MainWindowViewModel : ViewModelBase
    {
        private readonly ProxyService _proxyService;
        private string _selectedMode;
        private string _statusMessage;
        private string _statusColor;

        public MainWindowViewModel()
        {
            _proxyService = new ProxyService();
            Modes = new ObservableCollection<string> { "Off", "Global", "Pac" };
            _selectedMode = ""; // Start with no selection or Off? Python starts with -1 (no selection)
            _statusMessage = "Ready";
            _statusColor = "Black";
        }

        public ObservableCollection<string> Modes { get; }

        public string SelectedMode
        {
            get => _selectedMode;
            set
            {
                this.RaiseAndSetIfChanged(ref _selectedMode, value);
                OnModeChanged(value);
            }
        }

        public string StatusMessage
        {
            get => _statusMessage;
            set => this.RaiseAndSetIfChanged(ref _statusMessage, value);
        }

        public string StatusColor
        {
            get => _statusColor;
            set => this.RaiseAndSetIfChanged(ref _statusColor, value);
        }

        private void OnModeChanged(string mode)
        {
            if (string.IsNullOrEmpty(mode)) return;

            try
            {
                switch (mode)
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
}
