using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

public class ListUserMediaResponse
{
    [JsonPropertyName("images")]
    public List<MediaEntry> Images { set; get; }
}