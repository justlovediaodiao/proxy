using Avalonia.Controls;
using Avalonia.Controls.Templates;
using gui_net.ViewModels;
using gui_net.Views;

namespace gui_net;

/// <summary>
/// Given a view model, returns the corresponding view if possible.
/// </summary>
public class ViewLocator : IDataTemplate
{
    public Control? Build(object? param) => param switch
    {
        null => null,
        MainWindowViewModel => new MainWindow(),
        _ => new TextBlock { Text = $"Not Found: {param.GetType().FullName!.Replace("ViewModel", "View", StringComparison.Ordinal)}" }
    };

    public bool Match(object? data) => data is ViewModelBase;
}
