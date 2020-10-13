using System.Linq;
using System.Net;
using System.Threading.Tasks;
using Discord;
using Discord.Commands;
using TPUDISCORDBOT.Model;

namespace TPUDISCORDBOT.Modules
{
    public class CommandModule : ModuleBase<SocketCommandContext>
    {
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
                    client.DownloadFile(Context.Message.Attachments.ElementAt(0).Url, "./sounds/" + command + ".mp3");
                }
                var temp = new SoundModel();
                temp.command = command;
                temp.path = "./sounds/" + command + ".mp3";
                SoundManager.SoundManager.AddSound(temp);
                await Context.User.SendMessageAsync("Sound is uploaded with command " + command + ". Tell TPU to enable the command");
            }
            else
            {
                await Context.User.SendMessageAsync("Only .mp3 files are allowed. Please convert it first");
            }
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

    }
}