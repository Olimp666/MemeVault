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
        Console.WriteLine("Bot is running. Press Ctrl+C to exit.");
        
        // Ожидание бесконечно, пока не будет сигнала остановки
        await Task.Delay(Timeout.Infinite);
    }
}
