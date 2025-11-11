using System.Text.Json;
using MemeVaultControl.Model;

namespace MemeVaultControl.Client;

public class BackendClient
{
    private HttpClient _client = new();
    private string _serverUrl = "http://localhost";

    public async Task<UploadResponse?> UploadImage(UploadRequest uploadRequest)
    {
        var response = await _client.PostAsync(
            _serverUrl + "/upload",
            uploadRequest.ToForm()
        );

        response.EnsureSuccessStatusCode();
        var content = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<UploadResponse>(content);
    }

    public async Task<ListResponse?> GetList(string tag)
    {
        var response = await _client.GetAsync(
            _serverUrl + $"/images?tag={Uri.EscapeDataString(tag)}"
        );

        response.EnsureSuccessStatusCode();
        var content = await response.Content.ReadAsStringAsync();
        return JsonSerializer.Deserialize<ListResponse>(content);
    }
}