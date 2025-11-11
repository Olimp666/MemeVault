using MemeVaultControl.Commands;
using Telegram.Bot;
using Telegram.Bot.Exceptions;
using Telegram.Bot.Polling;
using Telegram.Bot.Types;

namespace MemeVaultControl.BotService;

public class ControlBotUpdateHandler : IUpdateHandler
{
    private readonly Dictionary<long, Command> _commands = [];

    public Task HandleUpdateAsync(ITelegramBotClient botClient, Update update, CancellationToken ct)
    {
        return update switch
        {
            { Message: { } msg } => HandleMessage(botClient, msg, ct),
            _ => Task.CompletedTask
        };
    }

    private async Task HandleMessage(ITelegramBotClient botClient, Message message, CancellationToken ct)
    {
        Console.WriteLine($"Handling message: [{message.Type}] {message.Text ?? message.Caption}");
        if (message.From is null) return;
        var userId = message.From.Id;
        var userHasCommand = _commands.TryGetValue(userId, out var command);

        if (!userHasCommand || command is null)
        {
            command = await ReadCommand(message, botClient, ct);
            if (command is null) return;
            _commands.Add(userId, command);
        }

        await command.Next(message);
        if (command.Finished) _commands.Remove(userId);
    }

    private async Task<Command?> ReadCommand(Message message, ITelegramBotClient botClient, CancellationToken ct)
    {
        var text = message.Text ?? message.Caption;

        if (text is null)
        {
            await botClient.SendMessage(
                message.Chat.Id,
                "Ошибка. Команда пустая",
                cancellationToken: ct
            );
            return null;
        }

        var cmd = text.Split(' ').FirstOrDefault()?.ToLower();

        if (cmd is null || !cmd.StartsWith('/'))
        {
            await botClient.SendMessage(
                message.Chat.Id,
                $"Ошибка. Команда должна начинаться с \"/\", но имеем {cmd}",
                cancellationToken: ct
            );
            return null;
        }

        Command? command = cmd switch
        {
            "/start" or "/help" => new StartCommand(botClient, ct),
            "/add" => new AddCommand(botClient, ct),
            "/list" => new ListCommand(botClient, ct),
            "/cancel" => new CancelCommand(botClient, ct),
            _ => null
        };

        if (command is null)
        {
            await botClient.SendMessage(
                message.Chat.Id,
                $"Ошибка. Неизвестная команда {cmd}",
                cancellationToken: ct
            );
            return null;
        }

        return command;
    }

    public Task HandleErrorAsync(ITelegramBotClient botClient, Exception exception, HandleErrorSource source,
        CancellationToken cancellationToken)
    {
        var errorMessage = exception switch
        {
            ApiRequestException apiRequestException
                => $"Telegram API Error:\n[{apiRequestException.ErrorCode}]\n{apiRequestException.Message}",
            _ => exception.ToString()
        };

        Console.WriteLine(errorMessage);
        return Task.CompletedTask;
    }
}