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

            var audioClient = await channel.ConnectAsync(); // Connect to channel
            await Say(audioClient, sound);                  // Play sound
            await channel.DisconnectAsync();                // Disconnect from channel

        }


        [Command("sounds")]
        [Alias("commands")]
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
            //await Context.User.SendMessageAsync("Please upload the audio file here. I will give further instructions once the file is downloaded on my end.");
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
