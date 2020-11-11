# TPU Discord bot

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cabdb7d6a21f4958b9a41ecdd2a67c4b)](https://app.codacy.com/gh/newtpuandre/TPUDISCORDBOT?utm_source=github.com&utm_medium=referral&utm_content=newtpuandre/TPUDISCORDBOT&utm_campaign=Badge_Grade)

A discord bot previously written in GO. Now rewritten in .net 5.0!

Users can add and play sounds in a discord voice channel. All sounds are moderated by an administrator

# Running with docker
* run docker build -t tpudiscordbot .

# Building yourself
Note: Make sure you have the .net 5.0 sdk installed.
* Git clone
* Run dotnet restore
* Run dotnet build
* Rename the app.config.sample file to app.config
* Insert your discord bot token.
* Run the compiled binary.