using System;
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
            string json = JsonConvert.SerializeObject(SoundManager.GetSounds().ToArray(), Formatting.Indented);
            File.WriteAllText(@soundlist, json);
        }

        public static List<SoundModel> readList()
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