using System.Net.Http.Headers;
using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class UploadRequest(long userId, string mediaId, IEnumerable<string> tags) : IForm
{
    public long UserId { get; set; } = userId;
    public string Image { get; set; } = mediaId;
    public IEnumerable<string> Tags { get; set; } = tags;

    public MultipartFormDataContent ToForm()
    {
        var form = new MultipartFormDataContent();

        form.Add(new StringContent(mediaId), "image");
        form.Add(new StringContent(JsonConvert.SerializeObject(Tags)), "tags");
        form.Add(new StringContent(UserId.ToString()), "user_id");

        return form;
    }
}