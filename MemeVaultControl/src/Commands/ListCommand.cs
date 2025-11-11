using MemeVaultControl.Client;
using MemeVaultControl.Model;
using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class ListCommand(ITelegramBotClient bot, CancellationToken ct) : CancellableCommand(bot, ct)
{
    private readonly BackendClient _client = new();
    private string[]? _tags;

    public override async Task Next(Message message)
    {
        await base.Next(message);
        if (Finished) return;

        TrySetTag(message);

        if (_tags is null)
        {
            await Reply(message, "По какому тегу искать?");
            return;
        }

        if (message.From is null) return;

        var images = await GetImages(message.From.Id, _tags);

        if (images is null)
        {
            await Reply(message, "При поиске произошел конфуз");
            return;
        }

        var formattedTags = string.Join(", ", _tags);

        if (images.Images.Count == 0)
        {
            await Reply(message, $"Для тега {formattedTags} нет совпадений");
            Finished = true;
            return;
        }

        const int bound = 10;

        var media = images.Images
            .Take(bound)
            .Select(CreateImage)
            .ToArray();

        await bot.SendMediaGroup(message.Chat.Id, media);
        var lessMessage = images.Images.Count < bound ? "" : "Показано {bound}";
        await Reply(message, $"Для тега {formattedTags} имеется {images.Images.Count} совпадений. {lessMessage}");

        Finished = true;
    }

    private InputMediaPhoto CreateImage(string fileId)
    {
        return new InputMediaPhoto(fileId);
    }

    private void TrySetTag(Message message)
    {
        var text = message.Text ?? message.Caption;
        if (text is null)
            return;

        var parts = text.Split(' ');
        var startIndex = text.StartsWith("/list") ? 1 : 0;

        if (parts.Length <= startIndex)
            return;

        _tags = parts[startIndex..];
    }

    private async Task<ListResponse?> GetImages(long userId, IEnumerable<string> tags)
    {
        try
        {
            var request = new ListRequest(userId, tags.ToList());
            var response = await _client.GetList(request);
            return response;
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }

        return null;
    }
}