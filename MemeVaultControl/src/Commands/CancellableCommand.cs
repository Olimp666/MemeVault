using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class CancellableCommand(ITelegramBotClient bot, CancellationToken ct) : Command(bot, ct)
{
    public override async Task Next(Message message)
    {
        var text = message.Text ?? message.Caption;

        if (text is null)
            return;

        if (!text.StartsWith("/cancel"))
            return;

        Finished = true;
        await Reply(message, "Команда отменена");
    }
}