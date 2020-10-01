using System;
using System.Threading.Tasks;
using Discord;
using Discord.WebSocket;
using System.Configuration;

namespace TPUDISCORDBOT
{
    class Program
    {

        private DiscordSocketClient _client;

        public static void Main(string[] args)
            => new Program().MainAsync().GetAwaiter().GetResult();

        public async Task MainAsync()
        {
            //Read config from app.config
            var appSettings = ConfigurationManager.AppSettings;
            var test = appSettings.Get("DiscordKey");
            Console.WriteLine("Read key : " + test);


            _client = new DiscordSocketClient();

            _client.Log += Log;

            await _client.LoginAsync(TokenType.Bot,
                Environment.GetEnvironmentVariable("DiscordToken"));
            await _client.StartAsync();

            // Block this task until the program is closed.
            await Task.Delay(-1);
        }

        private Task Log(LogMessage msg)
        {
            Console.WriteLine(msg.ToString());
            return Task.CompletedTask;
        }
    }
}
