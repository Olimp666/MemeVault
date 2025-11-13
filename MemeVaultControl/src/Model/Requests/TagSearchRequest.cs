using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class TagSearchRequest(long userId, List<string> tags)
{
    [JsonIgnore] public long UserId { get; set; } = userId;

    [JsonPropertyName("tags")] public List<string> Tags { get; set; } = tags;
}