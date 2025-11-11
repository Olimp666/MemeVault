using System.Text.Json.Serialization;

namespace MemeVaultInline.Models;

public class QueryResult
{
    public required string FileId { set; get; }
    public FileType Type { set; get; }
}

[JsonConverter(typeof(JsonStringEnumConverter))]
public enum FileType
{
    Photo,
    Gif,
    Video
}