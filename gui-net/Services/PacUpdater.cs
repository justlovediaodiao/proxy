using System.Net;
using System.Text;
using System.Text.Json;
using gui_net.Models;

namespace gui_net.Services;

public sealed class PacUpdater(ProcessLogBuffer logs)
{
    private const string GfwListUrl = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt";

    // Returns true when the PAC was rebuilt using the local fallback.
    public async Task<bool> UpdateAsync(Config config)
    {
        string? rules = null;
        try
        {
            rules = await DownloadAsync(null);
        }
        catch (Exception ex) when (IsDownloadFailure(ex))
        {
            logs.Append("pac", $"Direct download failed: {ex.Message}");
        }

        if (rules == null && config.Protocol is "http" or "socks5")
        {
            try
            {
                rules = await DownloadAsync(new WebProxy($"{config.Protocol}://{config.Host}:{config.Port}"));
            }
            catch (Exception ex) when (IsDownloadFailure(ex))
            {
                logs.Append("pac", $"Proxy download failed: {ex.Message}");
            }
        }

        var usedLocalRules = rules == null;
        if (usedLocalRules)
        {
            logs.Append("pac", "Online GFWList unavailable; using local rules.");
            rules = await File.ReadAllTextAsync("resources/gfwlist.txt");
        }

        var template = await File.ReadAllTextAsync("resources/abp.js");
        var userRules = await File.ReadAllTextAsync("resources/user-rule.txt");
        var content = template
            .Replace("__USERRULES__", SerializeRules(userRules))
            .Replace("__RULES__", SerializeRules(rules!))
            .Replace("__PROXY__", config.ProxyUrl);

        if (!usedLocalRules)
            await WriteAtomicAsync("resources/gfwlist.txt", rules!);
        await WriteAtomicAsync("resources/pac.js", content);
        logs.Append("pac", usedLocalRules ? "PAC rebuilt using local rules." : "PAC updated from online GFWList.");
        return usedLocalRules;
    }

    private static bool IsDownloadFailure(Exception ex) =>
        ex is HttpRequestException or TaskCanceledException or FormatException;

    private static async Task<string> DownloadAsync(IWebProxy? proxy)
    {
        using var handler = new HttpClientHandler { UseProxy = proxy != null, Proxy = proxy };
        using var client = new HttpClient(handler) { Timeout = TimeSpan.FromSeconds(15) };
        var encoded = await client.GetStringAsync(GfwListUrl);
        return Encoding.UTF8.GetString(Convert.FromBase64String(encoded));
    }

    private static string SerializeRules(string text)
    {
        var rules = text.Split('\n')
            .Select(line => line.TrimEnd('\r'))
            .Where(line => line.Length > 0 && !line.StartsWith('!') && !line.StartsWith('['))
            .ToArray();
        return JsonSerializer.Serialize(rules, JsonContext.Default.StringArray);
    }

    private static async Task WriteAtomicAsync(string path, string content)
    {
        var temporaryPath = path + "." + Guid.NewGuid().ToString("N") + ".tmp";
        try
        {
            await File.WriteAllTextAsync(temporaryPath, content);
            File.Move(temporaryPath, path, overwrite: true);
        }
        finally
        {
            if (File.Exists(temporaryPath))
                File.Delete(temporaryPath);
        }
    }
}
