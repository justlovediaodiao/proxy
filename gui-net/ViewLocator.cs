using System;

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
    public Control? Build(object? param)
    {
        if (param is null)
            return null;

        if (param is MainWindowViewModel)
        {
            return new MainWindow();
        }

        var name = param.GetType().FullName!.Replace("ViewModel", "View", StringComparison.Ordinal);
        return new TextBlock { Text = "Not Found: " + name };
    }

    public bool Match(object? data)
    {
        return data is ViewModelBase;
    }
}
