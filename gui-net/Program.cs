using Avalonia;
using Avalonia.ReactiveUI;
using gui_net;

BuildAvaloniaApp()
    .StartWithClassicDesktopLifetime(args);

// Avalonia configuration, don't remove; also used by visual designer.
static AppBuilder BuildAvaloniaApp()
    => AppBuilder.Configure<App>()
        .UsePlatformDetect()
        .WithInterFont()
        .LogToTrace()
        .UseReactiveUI();
