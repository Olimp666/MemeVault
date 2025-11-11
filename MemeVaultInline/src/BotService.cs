using MemeVaultInline.Helpers;
using MemeVaultInline.Models;
using Telegram.Bot.Exceptions;
using Telegram.Bot.Types;
using Telegram.Bot;
using Telegram.Bot.Polling;
using Telegram.Bot.Types.Enums;
using Telegram.Bot.Types.InlineQueryResults;
using RestSharp;
using FileType = MemeVaultInline.Models.FileType;

namespace MemeVaultInline;

internal class BotService
{
    private TelegramBotClient botClient;

    public BotService(TelegramBotClient botClient)
    {
        this.botClient = botClient;
        ReceiverOptions receiverOptions = new()
        {
            // AllowedUpdates = [UpdateType.InlineQuery, UpdateType.ChosenInlineResult]
            AllowedUpdates = []
        };
        botClient.StartReceiving(
            updateHandler: HandleUpdateAsync,
            errorHandler: HandleErrorAsync,
            receiverOptions: receiverOptions
        );
    }

    async Task HandleUpdateAsync(ITelegramBotClient _, Update update, CancellationToken cancellationToken)
    {
        switch (update.Type)
        {
            case UpdateType.InlineQuery:
                await OnInlineQueryReceived(update.InlineQuery!);
                break;
            case UpdateType.ChosenInlineResult:
                await OnChosenInlineResultReceived(update.ChosenInlineResult!);
                break;
        }
    }

    async Task OnInlineQueryReceived(InlineQuery query)
    {
        var results = new List<InlineQueryResult>();

        var client = new RestClient(ConfigHelper.Endpoint);
        var request = new RestRequest("/images", Method.Post);
        request.AddHeader("Content-Type", "application/json");
        var parsedTags = query.Query.Split().ToList();
        request.AddQueryParameter("user_id", query.From.Id);
        request.AddJsonBody(new { tags = parsedTags });
        var response = await client.ExecuteAsync<MatchResponse>(request);

        var memeList = response.Data?.ExactMatch;

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

    async Task OnChosenInlineResultReceived(ChosenInlineResult chosenInlineResult)
    {
        Console.WriteLine(chosenInlineResult.Query);
    }

    Task HandleErrorAsync(ITelegramBotClient bot, Exception exception, CancellationToken cancellationToken)
    {
        var errorMessage = exception switch
        {
            ApiRequestException apiRequestException
                => $"Ошибка Telegram API:\n[{apiRequestException.ErrorCode}]\n{apiRequestException.Message}",
            _ => exception.ToString()
        };

        Console.WriteLine(errorMessage);
        return Task.CompletedTask;
    }
}