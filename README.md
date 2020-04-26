# Twitch bot for tech streams

You need a `config.yml` file in the root of this repository to get the bot working

```yml
---
irc_token: xxx
nick: xxx
initial_channels:
    - xxx
prefix: "!"
```

And you should also create a `today.txt` file for your `!today` command content.

## Build & Deploy

I use Docker to build and run the bot. These are the commands I actually use.

## Build

```bash
docker build -t twitch-bot .
```

## Deploy

```bash
docker run -d --name twitch-bot -v /home/fabian/Dev/github.com/curiTTV/twitch-bot/today.txt:/opt/bot/today.txt --restart unless-stopped twitch-bot
```
