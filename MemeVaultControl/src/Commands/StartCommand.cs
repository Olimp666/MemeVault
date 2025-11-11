using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class StartCommand(ITelegramBotClient bot, CancellationToken ct) : Command(bot, ct)
{
    public override async Task Next(Message message)
    {
        await Reply(message, "Hi!\nTo add a meme, use /add\n");
        Finished = true;
    }
}