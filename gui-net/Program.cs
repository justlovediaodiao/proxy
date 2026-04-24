using Avalonia;
using gui_net;
using System.IO;
using System.Reflection;


Directory.SetCurrentDirectory(AppContext.BaseDirectory);

BuildAvaloniaApp()
    .StartWithClassicDesktopLifetime(args);

// Avalonia configuration, don't remove; also used by visual designer.
static AppBuilder BuildAvaloniaApp()
    => AppBuilder.Configure<App>()
        .UsePlatformDetect()
        .LogToTrace();
