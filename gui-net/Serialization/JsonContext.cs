using System.Text.Json.Serialization;
using gui_net.Models;

namespace gui_net.Serialization;

[JsonSourceGenerationOptions(
    PropertyNamingPolicy = JsonKnownNamingPolicy.CamelCase,
    WriteIndented = true
)]
[JsonSerializable(typeof(Config))]
public partial class JsonContext : JsonSerializerContext
{
}