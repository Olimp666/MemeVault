using System.Runtime.Serialization;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;

namespace MemeVaultControl.Model;

[JsonConverter(typeof(StringEnumConverter))]
public enum MediaType
{
    [EnumMember(Value = "photo")] Photo,
    [EnumMember(Value = "video")] Video,
    [EnumMember(Value = "gif")] Gif
}