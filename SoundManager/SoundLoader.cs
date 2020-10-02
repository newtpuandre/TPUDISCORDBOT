using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using Newtonsoft.Json;
using TPUDISCORDBOT.Model;

namespace TPUDISCORDBOT.SoundManager
{
    public static class SoundLoader
    {
        private static string soundlist = "./soundlist.json";

        public static void writeList()
        {
            var testModel = new SoundModel();
            testModel.path = "test";
            testModel.enabled = true;
            testModel.command = "hei";
            SoundManager.addSound(testModel);

            string json = JsonConvert.SerializeObject(SoundManager.GetSounds().ToArray(), Formatting.Indented);
            File.WriteAllText(@soundlist, json);
        }

        public static List<SoundModel> GetList()
        {
            using (StreamReader r = new StreamReader(soundlist))
            {
                string json = r.ReadToEnd();
                List<SoundModel> items = JsonConvert.DeserializeObject<List<SoundModel>>(json);
                return items;
            }
        }

    }
}