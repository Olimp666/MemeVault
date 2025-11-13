using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class IncreaseCounterRequest(long userId, string mediaId)
{
    [JsonIgnore] public long UserId { get; set; } = userId;
    [JsonIgnore] public string Image { get; set; } = mediaId;
}