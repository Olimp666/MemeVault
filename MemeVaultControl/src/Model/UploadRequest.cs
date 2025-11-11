using System.Net.Http.Headers;
using Newtonsoft.Json;

namespace MemeVaultControl.Model;

[JsonObject]
public class UploadRequest(long userId, byte[] image, IEnumerable<string> tags) : IForm
{
    public long UserId { get; set; } = userId;
    public byte[] Image { get; set; } = image;
    public IEnumerable<string> Tags { get; set; } = tags;

    public MultipartFormDataContent ToForm()
    {
        var form = new MultipartFormDataContent();

        var content = new ByteArrayContent(Image);
        content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");

        form.Add(content, "image", "content.jpg");
        form.Add(new StringContent(JsonConvert.SerializeObject(Tags)), "tags");
        form.Add(new StringContent(UserId.ToString()), "user_id");

        return form;
    }
}