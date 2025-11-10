using MemeVaultInline.Helpers;
class Program
{
    static async Task Main(string[] args)
    {
        Console.WriteLine($"Bot token: {ConfigHelper.BotToken}");
    }
}