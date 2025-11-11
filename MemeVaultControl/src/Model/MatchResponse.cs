using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class MediaItem
{
    [JsonPropertyName("tg_file_id")]
    public required string FileId { set; get; }
    
    [JsonPropertyName("file_type")]
    public FileType Type { set; get; }

    [JsonPropertyName("tags")]
    List<string> Tags { set; get; }
}

public class MediaList
{
    [JsonPropertyName("images")]
    public List<MediaItem> Images { set; get; }
}

public class MatchResponse
{
    [JsonPropertyName("exact_match")]
    public List<MediaItem> ExactMatch { get; set; }

    [JsonPropertyName("partial_match")]
    public List<MediaItem> PartialMatch { get; set; }
}

[JsonConverter(typeof(JsonStringEnumConverter))]
public enum FileType
{
    Photo,
    Gif,
    Video
}