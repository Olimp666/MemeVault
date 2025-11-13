using System.Configuration;

namespace MemeVaultControl.Helpers;

internal static class ConfigHelper
{
    public static readonly string BotToken = 
        Environment.GetEnvironmentVariable("Token") 
        ?? ConfigurationManager.AppSettings["Token"]!;
    
    public static readonly string Endpoint = 
        Environment.GetEnvironmentVariable("Endpoint") 
        ?? ConfigurationManager.AppSettings["Endpoint"]!;
}
