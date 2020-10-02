using System;
using System.Collections.Generic;
using TPUDISCORDBOT.Model;

namespace TPUDISCORDBOT.SoundManager
{
    public class SoundManager
    {

        private List<SoundModel> soundList;
        private SoundLoader _loader;
        public SoundManager()
        {
            soundList = new List<SoundModel>();
            _loader = new SoundLoader();
        }

    }
}