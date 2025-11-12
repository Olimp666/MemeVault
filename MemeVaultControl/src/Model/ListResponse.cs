using System.Text.Json.Serialization;
using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class ListResponse
{
    [JsonProperty("exact_match")] public required List<MediaEntry> ExactMatch { get; set; }
    [JsonProperty("partial_match")] public required List<MediaEntry> PartialMatch { get; set; }
}