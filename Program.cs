using System;
using System.Threading.Tasks;
using Discord;
using Discord.WebSocket;
using System.Configuration;
using Discord.Audio;
using System.Diagnostics;

namespace TPUDISCORDBOT
{
    class Program
    {

        private DiscordSocketClient _client;
        public DiscordSocketConfig _config;

        public static void Main(string[] args)
            => new Program().MainAsync().GetAwaiter().GetResult();

        public async Task MainAsync()
        {
            //Read config from app.config
            var appSettings = ConfigurationManager.AppSettings;
            var DiscordToken = appSettings.Get("DiscordKey");

            //Config DiscordSocketClient.
            _config = new DiscordSocketConfig { LogLevel = LogSeverity.Verbose };

            //New DiscordSocketClient with the config.
            _client = new DiscordSocketClient(_config);

            _client.Log += Logger.Log;
            _client.MessageReceived += MessageReceived;

            await _client.LoginAsync(TokenType.Bot, DiscordToken);
            await _client.StartAsync();

            // Block this task until the program is closed.
            await Task.Delay(-1);
        }

        private async Task MessageReceived(SocketMessage message)
        {
            if (message.Content == "!ping")
            {
                await message.Channel.SendMessageAsync("Pong!");
            }
        }


        //Edit object sound to a object containing info
        private async Task Say(IAudioClient connection, Object sound)
        {
            try
            {
                await connection.SetSpeakingAsync(true); // send a speaking indicator

                var psi = new ProcessStartInfo
                {
                    FileName = "ffmpeg",
                    Arguments = $@"-i ""{sound.Filename}"" -ac 2 -f s16le -ar 48000 pipe:1",
                    RedirectStandardOutput = true,
                    UseShellExecute = false
                };
                var ffmpeg = Process.Start(psi);

                var output = ffmpeg.StandardOutput.BaseStream;
                var discord = connection.CreatePCMStream(AudioApplication.Voice);
                await output.CopyToAsync(discord);
                await discord.FlushAsync();

                await connection.SetSpeakingAsync(false); // we're not speaking anymore
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
                Console.WriteLine($"- {ex.StackTrace}");
            }
        }

        //Edit object sound to a object containing info
        private async Task ConnectToVoice(SocketVoiceChannel voiceChannel)
        {
            if (voiceChannel == null)
                return;

            try
            {
                Console.WriteLine($"Connecting to channel {voiceChannel.Id}");
                var connection = await voiceChannel.ConnectAsync();
                Console.WriteLine($"Connected to channel {voiceChannel.Id}");


                await Task.Delay(1000);

                await Say(connection, Object.Hello);
            }
            catch (Exception ex)
            {
                // Oh no, error
                Console.WriteLine(ex.Message);
                Console.WriteLine($"- {ex.StackTrace}");
            }
        }
    }
}
