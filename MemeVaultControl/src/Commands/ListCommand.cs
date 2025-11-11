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

        var disposables = new List<IDisposable>();

        try
        {
            var media = images.Images
                .Take(bound)
                .Select(x => CreateImage(x, disposables))
                .ToArray();

            // Очень туго
            await bot.SendMediaGroup(message.Chat.Id, media);
            var lessMessage = images.Images.Count < bound ? "" : "Показано {bound}";
            await Reply(message, $"Для тега {_tag} имеется {images.Images.Count} совпадений. {lessMessage}");
        }
        finally
        {
            disposables.ForEach(x => x.Dispose());
        }

        Finished = true;
    }

    private InputMediaPhoto CreateImage(string base64, List<IDisposable> disposables)
    {
        // TODO: Here and in similar clauses add using
        var bytes = Convert.FromBase64String(base64);
        var stream = new MemoryStream(bytes);
        disposables.Add(stream);
        return new InputMediaPhoto(stream);
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