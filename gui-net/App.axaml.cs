using Avalonia;
using Avalonia.Controls.ApplicationLifetimes;
using Avalonia.Markup.Xaml;
using gui_net.Views;
using Avalonia.Media;
using Avalonia.Styling;
using SukiUI;
using SukiUI.Models;

namespace gui_net;

public partial class App : Application
{
    public override void Initialize()
    {
        AvaloniaXamlLoader.Load(this);
        ApplyColorTheme();
        ActualThemeVariantChanged += (_, _) => ApplyColorTheme();
    }

    private void ApplyColorTheme()
    {
        var dark = ActualThemeVariant == ThemeVariant.Dark;
        SukiTheme.GetInstance().ChangeColorTheme(new SukiColorTheme(
            "Cobalt", Color.Parse(dark ? "#9AAFFF" : "#315EEB"),
            Color.Parse(dark ? "#C0CDFF" : "#2449BC")));
    }

    public override void OnFrameworkInitializationCompleted()
    {
        if (ApplicationLifetime is IClassicDesktopStyleApplicationLifetime desktop)
        {
            desktop.MainWindow = new MainWindow();
        }

        base.OnFrameworkInitializationCompleted();
    }
}