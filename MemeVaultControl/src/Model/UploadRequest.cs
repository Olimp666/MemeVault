using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class UploadRequest(long userId, string mediaId, IEnumerable<string> tags)
{
    [JsonIgnore] public long UserId { get; set; } = userId;
    [JsonIgnore] public string Image { get; set; } = mediaId;
    [JsonProperty("tags")] public IEnumerable<string> Tags { get; set; } = tags;
}