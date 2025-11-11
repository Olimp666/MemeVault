using System.Text.Json.Serialization;
using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class ListResponse
{
    [JsonProperty("tg_file_ids")] public required List<string> Images { get; set; }
}