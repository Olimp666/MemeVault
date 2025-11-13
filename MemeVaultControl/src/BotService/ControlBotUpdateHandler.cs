using MemeVaultControl.Client;
using MemeVaultControl.Commands;
using Telegram.Bot;
using Telegram.Bot.Exceptions;
using Telegram.Bot.Polling;
using Telegram.Bot.Types;
using RestSharp;
using Telegram.Bot.Types.InlineQueryResults;
using MemeVaultControl.Model;
using MemeVaultControl.Helpers;

namespace MemeVaultControl.BotService;

public class ControlBotUpdateHandler : IUpdateHandler
{
    private readonly Dictionary<long, Command> _commands = [];
    private readonly BackendClient _client = new();

    public Task HandleUpdateAsync(ITelegramBotClient botClient, Update update, CancellationToken ct)
    {
        return update switch
        {
            { Message: { } msg } => HandleMessage(botClient, msg, ct),
            { InlineQuery: { } query } => HandleInlineQuery(botClient, query, ct),
            { ChosenInlineResult: { } result } => HandleChosenInlineResult(botClient, result, ct),
            _ => Task.CompletedTask
        };
    }

    async Task HandleChosenInlineResult(ITelegramBotClient bot, ChosenInlineResult result,
        CancellationToken ct)
    {
        var emptyQuery = result.Query.Length == 0;
        List<MediaEntry> memeList;
        if (!emptyQuery)
        {
            var request = new TagSearchRequest(result.From.Id, result.Query.Split().ToList());
            memeList = (await _client.SearchByTags(request))!.ExactMatch;
        }
        else
        {
            var request = new ListUserMediaRequest(result.From.Id);
            memeList = (await _client.ListUserMedia(request))!.Images;
        }

        var msgId = Convert.ToInt32(result.ResultId);
        await _client.IncreaseCounter(new IncreaseCounterRequest(result.From.Id, memeList[msgId].FileId));
        Console.WriteLine($"Incremented counter for image: [{memeList[msgId].FileId}]");
    }

    async Task HandleInlineQuery(ITelegramBotClient botClient, InlineQuery query, CancellationToken ct)
    {
        var emptyQuery = query.Query.Length == 0;
        var results = new List<InlineQueryResult>();
        List<MediaEntry> memeList;

        if (!emptyQuery)
        {
            var request = new TagSearchRequest(query.From.Id, query.Query.Split().ToList());
            memeList = (await _client.SearchByTags(request))!.ExactMatch;
        }
        else
        {
            var request = new ListUserMediaRequest(query.From.Id);
            memeList = (await _client.ListUserMedia(request))!.Images;
        }

        var cnt = 0;
        foreach (var item in memeList!)
        {
            InlineQueryResult result;
            switch (item.MediaType)
            {
                case MediaType.Photo:
                    result = new InlineQueryResultCachedPhoto(
                        id: cnt.ToString(),
                        photoFileId: item.FileId
                    );
                    break;
                case MediaType.Gif:
                    result = new InlineQueryResultCachedGif(
                        id: cnt.ToString(),
                        gifFileId: item.FileId
                    );
                    break;
                case MediaType.Video:
                    var vid = new InlineQueryResultCachedVideo(
                        id: cnt.ToString(),
                        videoFileId: item.FileId,
                        title: string.Join(" ", item.Tags) //TODO figure out video titles
                    );
                    vid.Caption = "Hi";
                    vid.Description = " ";
                    vid.ShowCaptionAboveMedia = true;
                    result = vid;
                    break;
                default:
                    return;
            }

            results.Add(result);
            cnt++;
        }

        var queryButton = new InlineQueryResultsButton("Загрузить мем")
        {
            StartParameter = "upload"
        };
        await botClient.AnswerInlineQuery(query.Id, results, cacheTime: 0, isPersonal: true,
            button: queryButton);
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

        var hasAttachment =
            message.Photo != null
            || message.Video != null
            || message.Animation != null;

        if (!hasAttachment && text is null)
        {
            await botClient.SendMessage(
                message.Chat.Id,
                "Ошибка. Команда пустая",
                cancellationToken: ct
            );
            return null;
        }

        text ??= "";
        var cmd = text.Split(' ').FirstOrDefault()?.ToLower();
        var args = text.Split(' ').Skip(1);

        if (!hasAttachment && (cmd is null || !cmd.StartsWith('/')))
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
            _ when hasAttachment => new AddCommand(botClient, ct),
            "/start" when args.FirstOrDefault() == "upload" => new AddHelpCommand(botClient, ct),
            "/start" or "/help" => new StartCommand(botClient, ct),
            "/add" when hasAttachment => new AddCommand(botClient, ct),
            "/add" => new AddHelpCommand(botClient, ct),
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