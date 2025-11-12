using System.Diagnostics;
using MemeVaultControl.Client;
using MemeVaultControl.Model;
using Telegram.Bot;
using Telegram.Bot.Types;
using Telegram.Bot.Types.Enums;

namespace MemeVaultControl.Commands;

public class AddCommand(ITelegramBotClient bot, CancellationToken ct) : CancellableCommand(bot, ct)
{
    private string? _media;
    private MediaType? _mediaType;
    private string[]? _tags;
    private readonly BackendClient _client = new();

    public override async Task Next(Message message)
    {
        await base.Next(message);
        if (Finished) return;

        if (message.From?.Id is null) return;

        TrySetMedia(message);
        TrySetTags(message);

        if (_media is null || _mediaType is null)
        {
            await Reply(message, "Приложите фотографию");
            return;
        }

        if (_tags is null)
        {
            await Reply(message, "Предоставьте теги");
            return;
        }

        Finished = true;
        var success = await SendRequest(message.From!.Id, _media, _mediaType.Value, _tags);

        if (!success)
        {
            await Reply(message, "Произошел конфуз при сохранении мема");
            return;
        }

        await SignalSuccess(message, _tags);
    }

    private void TrySetMedia(Message message)
    {
        if (_media is not null) return;

        _mediaType = message.Type switch
        {
            MessageType.Photo => MediaType.Photo,
            MessageType.Video => MediaType.Video,
            MessageType.Animation => MediaType.Gif,
            _ => null
        };

        _media = message.Photo?.LastOrDefault()?.FileId
                 ?? message.Video?.FileId
                 ?? message.Animation?.FileId;
    }

    private void TrySetTags(Message message)
    {
        if (_tags is not null) return;

        var text = message.Text ?? message.Caption;
        if (text is null)
            return;

        var parts = text.Split(' ');
        var startIndex =
            text.StartsWith("/add") ? 1 :
            text.StartsWith("/start") ? 2 : 0;

        if (parts.Length <= startIndex)
            return;

        _tags = parts[startIndex..];
    }

    private async Task SignalSuccess(Message message, string[] tags)
    {
        var formattedTags = string.Join(", ", tags);
        await Reply(message, $"Мем успешно сохранен с тегами [{formattedTags}]");
    }

    private async Task<bool> SendRequest(long userId, string media, MediaType mediaType, string[] tags)
    {
        var form = new UploadRequest(userId, media, mediaType, tags);
        try
        {
            await _client.UploadImage(form);
            return true;
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }

        return false;
    }
}