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

    public Task HandleUpdateAsync(ITelegramBotClient botClient, Update update, CancellationToken ct)
    {
        return update switch
        {
            { Message: { } msg } => HandleMessage(botClient, msg, ct),
            { InlineQuery: { } query } => HandleInlineQuery(botClient, query, ct),
            _ => Task.CompletedTask
        };
    }

    async Task HandleInlineQuery(ITelegramBotClient botClient, InlineQuery query, CancellationToken ct)
    {
        var emptyQuery = query.Query.Length == 0;
        var results = new List<InlineQueryResult>();
        List<MediaItem> memeList;

        var client = new RestClient(ConfigHelper.Endpoint);


        if (!emptyQuery)
        {
            var request = new RestRequest("/images", Method.Post);
            request.AddQueryParameter("user_id", query.From.Id);
            var parsedTags = query.Query.Split().ToList();
            request.AddJsonBody(new { tags = parsedTags });
            var response = await client.ExecuteAsync<MatchResponse>(request);
            memeList = response.Data?.ExactMatch!;
        }
        else
        {
            var request = new RestRequest("/user/images");
            request.AddQueryParameter("user_id", query.From.Id);
            var response = await client.ExecuteAsync<MediaList>(request);
            memeList = response.Data!.Images;
        }

        var cnt = 0;
        foreach (var item in memeList!)
        {
            InlineQueryResult result;
            switch (item.Type)
            {
                case FileType.Photo:
                    result = new InlineQueryResultCachedPhoto(
                        id: cnt.ToString(),
                        photoFileId: item.FileId
                    );
                    break;
                case FileType.Gif:
                    result = new InlineQueryResultCachedGif(
                        id: cnt.ToString(),
                        gifFileId: item.FileId
                    );
                    break;
                case FileType.Video:
                    result = new InlineQueryResultCachedVideo(
                        id: cnt.ToString(),
                        videoFileId: item.FileId,
                        title: "⠀" //TODO figure out video titles
                    );
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