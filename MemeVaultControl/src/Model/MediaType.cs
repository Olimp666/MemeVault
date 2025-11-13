using System.Runtime.Serialization;
using System.Text.Json.Serialization;

namespace MemeVaultControl.Model;

[JsonConverter(typeof(JsonStringEnumConverter))]
public enum MediaType
{
    Photo,
    Video,
    Gif
}