using System;
using System.Collections.Generic;
using TPUDISCORDBOT.Model;

namespace TPUDISCORDBOT.SoundManager
{
    public static class SoundManager
    {

        private static List<SoundModel> soundList = new List<SoundModel>();

        public static List<SoundModel> GetSounds()
        {
            return soundList;
        }

        public static void AddSound(SoundModel item)
        {
            soundList.Add(item);
            SoundLoader.writeList();
        }

        public static void UpdateSound(SoundModel item)
        {

        }

        public static SoundModel GetSound(string command)
        {
            return soundList.Find(x => x.command.ToLower() == command.ToLower());
        }

        public static List<SoundModel> GetList()
        {
            return soundList;
        }
        public static void SetList(List<SoundModel> list)
        {
            soundList = list;
        }

    }
}