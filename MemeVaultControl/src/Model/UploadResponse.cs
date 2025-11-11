using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

[Serializable]
public class UploadResponse
{
    [JsonPropertyName("image_id")]
    public long ImageId { get; set; }
}