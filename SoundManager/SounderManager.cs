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

        public static void addSound(SoundModel item)
        {
            soundList.Add(item);
        }

    }
}