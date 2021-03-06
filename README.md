# discord-bot

## Getting started with development
You'll need a discord bot if you want to do any manual testing, to do that I'd recommend creating your own discord server to use.
Then visit this page https://discord.com/developers/applications and create a new bot.
Then follow this tutorial (its short) https://discordjs.guide/preparations/adding-your-bot-to-servers.html
The bot will have a token that you can copy paste, you use that whenever you run the bot program.

The bot uses a MongoDB database, so go install that.

Now install `go`, you'll need at least version 1.16.
A couple things you'll need for development
```
go build main.go     <-- compiles your code
go run main.go       <-- compiles and runs your code
```

The bot is configured using the `.bot.conf.json` file, a basic config would look like this:
```
{
    "token" : "my bots token",
    "mongoDatabase" : {
        "dbname": "discordbots",
        "quoteCollection": "quotes",
        "connStr": "mongodb://localhost:27017"
    },
    "settings" : {
        "quotechannel" : "837625133635993602",
        "leavemessagechannel" : "837707405626048552"
    },
    "enabledFeatures" : [
        "leavemessage"
    ],
    "enabledCommands" : [
        "help",
        "getquote",
        "addquote"
    ],
    "commandPrefix" : "$"
}
```