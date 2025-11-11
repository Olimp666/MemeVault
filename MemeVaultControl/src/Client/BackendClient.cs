using MemeVaultControl.Helpers;
using MemeVaultControl.Model;
using Newtonsoft.Json;

namespace MemeVaultControl.Client;

public class BackendClient
{
    private readonly HttpClient _client = new();
    private readonly string _serverUrl = ConfigHelper.ServerUrl;

    public async Task UploadImage(UploadRequest uploadRequest)
    {
        var body = new StringContent(JsonConvert.SerializeObject(uploadRequest));
        var response = await _client.PostAsync(
            _serverUrl + $"/upload?user_id={uploadRequest.UserId}&tg_file_id={uploadRequest.Image}",
            body
        );
        
        response.EnsureSuccessStatusCode();
    }

    public async Task<ListResponse?> GetList(ListRequest listRequest)
    {
        var body = new StringContent(JsonConvert.SerializeObject(listRequest));
        var response = await _client.PostAsync(
            _serverUrl + $"/images?user_id={listRequest.UserId}",
            body
        );

        var content = await response.Content.ReadAsStringAsync();
        response.EnsureSuccessStatusCode();
        return  JsonConvert.DeserializeObject<ListResponse>(content);
    }
}