using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class StartCommand(ITelegramBotClient bot, CancellationToken ct) : Command(bot, ct)
{
    public override async Task Next(Message message)
    {
        await Reply(message, "Привет!\nЧтобы добавить мем, приложите фотографию/видео/гифку вместе с тегами сообщением боту\n");
        Finished = true;
    }
}