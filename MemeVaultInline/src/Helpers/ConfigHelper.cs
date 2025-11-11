using System.Configuration;

namespace MemeVaultInline.Helpers;

internal static class ConfigHelper
{
    public static readonly string BotToken = ConfigurationManager.AppSettings["Token"]!;
    public static readonly string Endpoint = ConfigurationManager.AppSettings["Endpoint"]!;
}