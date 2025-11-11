using System.Diagnostics;
using MemeVaultControl.Client;
using MemeVaultControl.Model;
using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class AddCommand(ITelegramBotClient bot, CancellationToken ct) : CancellableCommand(bot, ct)
{
    private string? _media;
    private string[]? _tags;
    private readonly BackendClient _client = new();

    public override async Task Next(Message message)
    {
        await base.Next(message);
        if (Finished) return;

        await TrySetMedia(message);
        TrySetTags(message);

        if (_media is null)
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
        Debug.Assert(message.From?.Id is not null);
        var success = await SendRequest(message.From!.Id, _media, _tags);

        if (!success)
        {
            await Reply(message, "Произошел конфуз при сохранении мема");
            return;
        }

        await SignalSuccess(message, _tags);
    }

    private async Task TrySetMedia(Message message)
    {
        var photos = message.Photo;

        if (photos is null)
            return;

        if (photos.Length == 0)
            return;

        // TODO: Pick desirable quality
        var photo = photos.Last();

        _media = photo.FileId;
    }

    private void TrySetTags(Message message)
    {
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

    private async Task<bool> SendRequest(long userId, string media, string[] tags)
    {
        var form = new UploadRequest(userId, media, tags);
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