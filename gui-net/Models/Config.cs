using System.Text.Json.Serialization;

namespace gui_net.Models;

[JsonSourceGenerationOptions(
    PropertyNamingPolicy = JsonKnownNamingPolicy.CamelCase,
    WriteIndented = true
)]
[JsonSerializable(typeof(Config))]
public partial class JsonContext : JsonSerializerContext;

public class Config
{
    [JsonPropertyName("host")]
    public string Host { get; set; } = "127.0.0.1";

    [JsonPropertyName("port")]
    public int Port { get; set; } = 1080;

    [JsonPropertyName("protocol")]
    public string Protocol { get; set; } = "socks5";

    [JsonPropertyName("pac_host")]
    public string PacHost { get; set; } = "127.0.0.1";

    [JsonPropertyName("pac_port")]
    public int PacPort { get; set; } = 1080;

    [JsonPropertyName("proxy_commands")]
    public string[] ProxyCommands { get; set; } = [];

    [JsonIgnore]
    public string ProxyUrl { get; set; } = "";
}
