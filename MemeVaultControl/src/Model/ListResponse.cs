using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

[Serializable]
public class ListResponse
{
    [JsonPropertyName("images")]
    public required List<string> Images { get; set; }
}