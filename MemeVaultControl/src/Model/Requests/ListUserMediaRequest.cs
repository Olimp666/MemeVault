using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class ListUserMediaRequest(long userId)
{
    [JsonIgnore] public long UserId { get; set; } = userId;
}