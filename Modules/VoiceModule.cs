using System;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Threading.Tasks;
using Discord;
using Discord.Audio;
using Discord.Commands;
using TPUDISCORDBOT.Model;
using TPUDISCORDBOT.Services;
using TPUDISCORDBOT.SoundManager;

namespace TPUDISCORDBOT.Modules
{

    // Modules must be public and inherit from an IModuleBase
    public class VoiceModule : ModuleBase<SocketCommandContext>
    {

        [Command("play", RunMode = RunMode.Async)]
        public async Task JoinChannel(string command = null, IVoiceChannel channel = null)
        {
            if (Program.playingSound)
            {
                await Context.User.SendMessageAsync("I am currently busy. Try again shortly");
                return;
            }
            // Get the audio channel
            channel = channel ?? (Context.User as IGuildUser)?.VoiceChannel;
            if (channel == null) { await Context.Channel.SendMessageAsync("User must be in a voice channel, or a voice channel must be passed as an argument."); return; }
            if (command == null) { await Context.User.SendMessageAsync("Please provide a sound name. Example: !play name. !sounds for a full list of available sounds."); return; }

            var sound = SoundManager.SoundManager.GetSound(command);
            if (sound == null)
            {
                await Context.User.SendMessageAsync("The sound you have requested does not exist. Please try again");
                return;
            }

            if (!sound.enabled)
            {
                await Context.User.SendMessageAsync("The sound you have requested is not enabled. Please tell TPU to activate it before trying again");
                return;
            }
            Program.playingSound = true;
            var audioClient = await channel.ConnectAsync(); // Connect to channel
            Program.audioClient = audioClient;
            await Say(audioClient, sound);                  // Play sound
            await channel.DisconnectAsync();                // Disconnect from channel
            Program.playingSound = false;
        }

        [Command("stop")]
        public async Task Stop()
        {
            await Program.audioClient.SetSpeakingAsync(false);
            await Program.audioClient.StopAsync();
        }

        [Command("sounds")]
        public async Task GetSoundList()
        {
            EmbedBuilder builder = new EmbedBuilder();
            builder.WithTitle("TPUBOT Sound list");
            builder.WithDescription("To play sound : !play <name>. Replace <name> with one of the names below.");
            foreach (var item in SoundManager.SoundManager.GetList())
            {
                builder.AddField("Sound name", item.command);
            }

            await Context.Channel.SendMessageAsync("", false, builder.Build());
        }

        [Command("commands")]
        [Alias("command")]
        public async Task GetCommands()
        {
            EmbedBuilder builder = new EmbedBuilder();
            builder.WithTitle("TPUBOT Command List");
            builder.AddField("!commands", "!commands");
            builder.AddField("!sounds", "!sounds");
            builder.AddField("!upload <name>", "!upload <name>");
            builder.AddField("!play <name>", "!play <name>");
            builder.AddField("!stop", "!stop");

            await Context.Channel.SendMessageAsync("", false, builder.Build());
        }

        [Command("toggle")]
        public async Task ToggleSound(string command = null)
        {
            //Check username
            if (command != null && Context.User.Id == 132626622816845824)
            {

                var result = SoundManager.SoundManager.ToggleSound(command);
                await Context.User.SendMessageAsync($"{command} state is now: {result}");

            }
            else
            {
                await Context.User.SendMessageAsync("Only TPU is able to enable commands");
            }
        }

        [Command("upload", RunMode = RunMode.Async)]
        public async Task UploadSound(string command = null)
        {
            if (command == null && Context.Message.Attachments.Count <= 0)
            {
                await Context.User.SendMessageAsync("Please upload the file here and use the following command. !upload soundName");
                return;
            }
            if (Context.Message.Attachments.ElementAt(0).Filename.Contains(".mp3"))
            {
                using (var client = new WebClient())
                {
                    client.DownloadFile(Context.Message.Attachments.ElementAt(0).Url, "./" + command + ".mp3");
                }
                var temp = new SoundModel();
                temp.command = command;
                temp.path = "./" + command + ".mp3";
                SoundManager.SoundManager.AddSound(temp);
                await Context.User.SendMessageAsync("Sound is uploaded with command " + command + ". Tell TPU to enable the command");
            }
            else
            {
                await Context.User.SendMessageAsync("Only .mp3 files are allowed. Please convert it first");
                return;
            }
        }

        private static async Task Say(IAudioClient connection, SoundModel sound)
        {
            try
            {
                await connection.SetSpeakingAsync(true); // send a speaking indicator

                var psi = new ProcessStartInfo
                {
                    FileName = "ffmpeg",
                    Arguments = $@"-i ""{sound.path}"" -ac 2 -f s16le -ar 48000 pipe:1",
                    RedirectStandardOutput = true,
                    UseShellExecute = false
                };
                var ffmpeg = Process.Start(psi);

                var output = ffmpeg.StandardOutput.BaseStream;
                var discord = connection.CreatePCMStream(AudioApplication.Mixed);
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

    }
}
