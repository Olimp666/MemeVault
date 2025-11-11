using Telegram.Bot;
using Telegram.Bot.Polling;
using Telegram.Bot.Types.Enums;

namespace MemeVaultControl.BotService;

public class BotService(ITelegramBotClient bot)
{
    private readonly CancellationTokenSource _cts = new();

    public Task Run()
    {
        ReceiverOptions receiverOptions = new()
        {
            AllowedUpdates = [UpdateType.Message] // receive all update types except ChatMember related updates
        };

        bot.StartReceiving(
            updateHandler: new ControlBotUpdateHandler(),
            receiverOptions: receiverOptions,
            cancellationToken: _cts.Token
        );

        return Task.CompletedTask;
    }
}