using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public abstract class Command(ITelegramBotClient bot, CancellationToken ct)
{
    public bool Finished = false;
    public abstract Task Next(Message message);

    protected async Task Reply(Message message, string reply)
        => await bot.SendMessage(message.Chat.Id, reply);
}