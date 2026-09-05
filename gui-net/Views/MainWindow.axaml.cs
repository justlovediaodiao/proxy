using Avalonia.Controls;
using Avalonia.Input;
using gui_net.ViewModels;

namespace gui_net.Views;

public partial class MainWindow : Window
{
    private LogWindow? _logWindow;

    public MainWindow()
    {
        InitializeComponent();
        DataContext = new MainWindowViewModel();
        Closing += MainWindow_Closing;
    }

    private void ProxyTitle_PointerPressed(object? sender, PointerPressedEventArgs e)
    {
        if (DataContext is not MainWindowViewModel vm)
            return;

        if (_logWindow != null)
        {
            _logWindow.Activate();
            return;
        }

        _logWindow = new LogWindow(vm.Logs);
        _logWindow.Closed += (_, _) => _logWindow = null;
        _logWindow.Show(this);
        e.Handled = true;
    }

    private void MainWindow_Closing(object? sender, System.ComponentModel.CancelEventArgs e)
    {
        if (DataContext is MainWindowViewModel vm)
            vm.OnExit();
    }
}
