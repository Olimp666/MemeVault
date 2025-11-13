using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class MediaEntry
{
    [JsonPropertyName("tg_file_id")] public string FileId { get; set; }
    [JsonPropertyName("file_type")] public MediaType MediaType { get; set; }
    [JsonPropertyName("tags")] public List<string> Tags { get; set; }
}