using MemeVaultInline.Helpers;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.InlineQueryResults;

namespace MemeVaultInline;

internal static class Program
{
    private static Task Main(string[] args)
    {
        var botClient =
            new TelegramBotClient(ConfigHelper.BotToken, cancellationToken: new CancellationTokenSource().Token);
        var botService = new BotService(botClient);
        Console.WriteLine($"Bot is running...");
        Console.ReadLine();
        return Task.CompletedTask;
    }
}