using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class MediaEntry
{
    [JsonProperty("tg_file_id")] public string FileId { get; set; }
    [JsonProperty("file_type")] public MediaType MediaType { get; set; }
}