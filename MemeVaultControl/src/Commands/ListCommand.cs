using MemeVaultControl.Client;
using MemeVaultControl.Model;
using Telegram.Bot;
using Telegram.Bot.Types;

namespace MemeVaultControl.Commands;

public class ListCommand(ITelegramBotClient bot, CancellationToken ct) : CancellableCommand(bot, ct)
{
    private readonly BackendClient _client = new();
    private string? _tag;

    public override async Task Next(Message message)
    {
        await base.Next(message);
        if (Finished) return;

        TrySetTag(message);

        if (_tag is null)
        {
            await Reply(message, "По какому тегу искать?");
            return;
        }

        var images = await GetImages(_tag);

        if (images is null)
        {
            await Reply(message, "При поиске произошел конфуз");
            return;
        }

        const int bound = 10;
        
        var media = images.Images
            .Take(bound)
            .Select(CreateImage)
            .ToArray();

        // Очень туго
        await bot.SendMediaGroup(message.Chat.Id, media);
        var lessMessage = images.Images.Count < bound ? "" : "Показано {bound}";
        await Reply(message, $"Для тега {_tag} имеется {images.Images.Count} совпадений. {lessMessage}");

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

        _tag = parts[startIndex];
    }

    private async Task<ListResponse?> GetImages(string tag)
    {
        try
        {
            var response = await _client.GetList(tag);
            return response;
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }

        return null;
    }
}