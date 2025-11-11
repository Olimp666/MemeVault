using MemeVaultControl.Helpers;
using Telegram.Bot;

namespace MemeVaultControl;

public static class Program
{
    public static async Task Main(string[] args)
    {
        var bot = new TelegramBotClient(ConfigHelper.BotToken);
        var botService = new BotService.BotService(bot);
        await botService.Run();
        Console.WriteLine("Bot is running.\nPress any key to exit...");
        Console.ReadKey();
    }
}