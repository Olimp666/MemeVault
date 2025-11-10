using System.Configuration;

namespace MemeVaultInline.Helpers;

internal static class ConfigHelper
{
    public static string BotToken = ConfigurationManager.AppSettings["Token"]!;
}