using System.Threading;
using Avalonia.Controls;
using Avalonia.Threading;
using Avalonia.VisualTree;
using gui_net.Services;

namespace gui_net.Views;

public partial class LogWindow : Window
{
    private readonly ProcessLogBuffer _logs;
    private readonly DispatcherTimer _refreshTimer;
    private int _hasChanges = 1;

    public LogWindow() : this(new ProcessLogBuffer())
    {
    }

    public LogWindow(ProcessLogBuffer logs)
    {
        InitializeComponent();
        _logs = logs;
        _logs.Changed += Logs_Changed;

        _refreshTimer = new DispatcherTimer
        {
            Interval = TimeSpan.FromMilliseconds(150)
        };
        _refreshTimer.Tick += RefreshTimer_Tick;
        _refreshTimer.Start();

        Closed += LogWindow_Closed;
        RefreshLog();
    }

    private void Logs_Changed()
    {
        Interlocked.Exchange(ref _hasChanges, 1);
    }

    private void RefreshTimer_Tick(object? sender, EventArgs e)
    {
        if (Interlocked.Exchange(ref _hasChanges, 0) == 1)
            RefreshLog();
    }

    private void RefreshLog()
    {
        LogTextBox.Text = _logs.GetSnapshot();
        Dispatcher.UIThread.Post(() =>
        {
            LogTextBox.GetVisualDescendants()
                .OfType<ScrollViewer>()
                .FirstOrDefault()?
                .ScrollToEnd();
        }, DispatcherPriority.Background);
    }

    private void LogWindow_Closed(object? sender, EventArgs e)
    {
        _refreshTimer.Stop();
        _logs.Changed -= Logs_Changed;
    }
}
