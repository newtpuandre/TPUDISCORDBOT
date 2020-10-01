using System;
using System.Threading.Tasks;
using Discord;

/**
 TODO:
 [] Log to file
 [] Log with filestamp
 [] Log errors
 [] Log based on the severity level set in main

*/


namespace TPUDISCORDBOT
{
    public static class Logger
    {
        public static Task Log(LogMessage msg)
        {
            Console.WriteLine(msg.ToString());
            return Task.CompletedTask;
        }
    }

}
