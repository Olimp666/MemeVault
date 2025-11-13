using MemeVaultControl.Helpers;
using MemeVaultControl.Model;
using Newtonsoft.Json;
using RestSharp;

namespace MemeVaultControl.Client;

public class BackendClient
{
    private readonly RestClient _client = new(ConfigHelper.Endpoint);

    public async Task UploadImage(UploadRequest body)
    {
        var request = new RestRequest("/upload", Method.Post);
        request.AddQueryParameter("user_id", body.UserId);
        request.AddQueryParameter("tg_file_id", body.Image);
        request.AddQueryParameter("file_type", body.MediaType.ToString().ToLower());
        request.AddJsonBody(new { tags = body.Tags });
        var response = await _client.ExecuteAsync(request);
        
        Console.WriteLine($"Executed /upload with status code {response.StatusCode}");
    }

    public async Task<TagSearchResponse?> SearchByTags(TagSearchRequest body)
    {
        var request = new RestRequest("/images", Method.Post);
        request.AddQueryParameter("user_id", body.UserId);
        var parsedTags = body.Tags;
        request.AddJsonBody(new { tags = parsedTags });
        var response = await _client.ExecuteAsync<TagSearchResponse>(request);
        return response.Data;
    }

    public async Task<ListUserMediaResponse?> ListUserMedia(ListUserMediaRequest body)
    {
        var request = new RestRequest("/user/images");
        request.AddQueryParameter("user_id", body.UserId);
        var response = await _client.ExecuteAsync<ListUserMediaResponse>(request);
        return response.Data;
    }
}