using System.Text.Json.Serialization;

namespace MemeVaultInline.Models;

public class MatchItem
{
    [JsonPropertyName("tg_file_id")]
    public required string FileId { set; get; }
    
    [JsonPropertyName("file_type")]
    public FileType Type { set; get; }
}

public class MatchResponse
{
    [JsonPropertyName("exact_match")]
    public List<MatchItem> ExactMatch { get; set; }

    [JsonPropertyName("partial_match")]
    public List<MatchItem> PartialMatch { get; set; }
}

[JsonConverter(typeof(JsonStringEnumConverter))]
public enum FileType
{
    Photo,
    Gif,
    Video
}