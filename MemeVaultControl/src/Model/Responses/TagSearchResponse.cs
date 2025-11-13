using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;
public class TagSearchResponse
{
    [JsonPropertyName("exact_match")] public required List<MediaEntry> ExactMatch { get; set; }
    [JsonPropertyName("partial_match")] public required List<MediaEntry> PartialMatch { get; set; }
}