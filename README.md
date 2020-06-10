# Twitch bot for twitch.tv/curi

The old bot was written in Python and can be found in the [python branch](https://github.com/curiTTV/twitch-bot/tree/python).

Two env vars have to be set `TOKEN` (Twitch chat auth token) and `MONGO_HOST` (MongoDB hostname or ip). For the connection to the MongoDB host the default port 27017 will be used.

## MongoDB

The database `bots` and two collections named `commands` and `talers` will be created after the usage of their respective built-in commands.
