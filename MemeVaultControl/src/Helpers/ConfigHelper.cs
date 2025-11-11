using System.Configuration;

namespace MemeVaultControl.Helpers;

internal static class ConfigHelper
{
    public static readonly string BotToken = ConfigurationManager.AppSettings["Token"]!;
}