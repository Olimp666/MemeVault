using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class ListRequest(long userId, List<string> tags)
{
    [JsonIgnore] public long UserId { get; set; } = userId;

    [JsonProperty("tags")] public List<string> Tags { get; set; } = tags;
}