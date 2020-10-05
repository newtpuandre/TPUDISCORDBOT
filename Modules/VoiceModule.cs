using System;
using System.Collections.Generic;
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
            if (Program.audioClient != null)
            {
                await Program.audioClient.SetSpeakingAsync(false);
                await Program.audioClient.StopAsync();
                Program.audioClient = null;
            }
        }

        public static byte[][] ReadOpus(SoundModel sound)
        {
            // Int16 opuslen 2bytes
            //byte[][] temp;
            List<byte[]> byteList = new List<byte[]>();

            BinaryReader binReader = new BinaryReader(File.Open(sound.path, FileMode.Open));
            try
            {
                while (binReader.BaseStream.Position != binReader.BaseStream.Length)
                {
                    Console.WriteLine("Reading bytes");
                    var readBytes = binReader.ReadBytes(4);
                    byteList.Add(readBytes);
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex);
            }
            binReader.Close();
            return byteList.ToArray();
        }

        private static async Task Say(IAudioClient connection, SoundModel sound)
        {
            try
            {
                await connection.SetSpeakingAsync(true); // send a speaking indicator
                if (sound.path.Contains(".mp3"))
                {
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

                }
                else
                { //Bad DCA section
                    byte[][] dcaFile = ReadOpus(sound);
                    Console.WriteLine("DCAFILE Length " + dcaFile.Length);
                    var discord = connection.CreatePCMStream(AudioApplication.Mixed);
                    int i = 0;
                    MemoryStream stream = new MemoryStream();
                    foreach (var item in dcaFile)
                    {
                        await stream.WriteAsync(item);
                        Console.WriteLine("Write new data . " + i++ + " Item length. " + item.Length);
                    }
                    await stream.CopyToAsync(discord);
                    await discord.FlushAsync();
                    //discord.Flush();

                }

                await connection.SetSpeakingAsync(false); // we're not speaking anymore
                Program.playingSound = false;
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
                Console.WriteLine($"- {ex.StackTrace}");
            }
        }

    }
}
