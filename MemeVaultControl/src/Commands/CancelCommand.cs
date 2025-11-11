using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class CancelCommand(ITelegramBotClient bot, CancellationToken ct) : Command(bot, ct)
{
    public override async Task Next(Message message)
    {
        await Reply(message, "Нечего отменять");
        Finished = true;
    }
}