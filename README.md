# Pinbot

[Install 📌](https://discord.com/oauth2/authorize?client_id=921554139740254209&permissions=3136&redirect_uri=https%3A%2F%2Fgithub.com%2Felliotwms%2Fpinbot&scope=applications.commands%20bot)

[Join the Pinbot Discord Server](https://discord.gg/a3u2PZ6V28)

Whenever you use the bot's "Pin" command (right-click a message, choose Apps, choose Pin), Pinbot posts the message to a channel.

Pinbot will reply with a link to the pinned message, and signal it's done by reacting to the original message with a 📌 emoji.

![Example of a Pinbot message](https://user-images.githubusercontent.com/4396779/147515477-850ab41a-6a89-4746-9f65-e27c259f7602.png)

### Why does this exist?

Pinbot is designed as an extension to Discord's channel pins system. Use Pinbot to:
* Bypass Discord's 50-pin limit and create a historic stream of all your pins
* Collect all your server's pins into one place (with optional overrides)
* Give your server's pins a more permanent home

Discord guilds use pins for a lot more than just highlighting important information. In fact, many guilds use the pin system as a form of memorialising a good joke, a savage putdown, or other memorable moments. As a result, the 50 pin per channel limit means that in order to keep something, you will eventually have to get rid of something else.

### How does it work?

Pinbot uses the channel name to decide where it will post. In order of priority it will pin in:
1. `#{channel}-pins`, where `channel` is the name of the channel the message was pinned in
2. `#pins`, a general pins channel
3. `#{channel}`, the channel the pin was posted in, so that if you don't want a separate pins channel you can instead 
search for pins by @pinbot in the channel

Whenever Pinbot pins a message, or whenever you update the actual channel pins, Pinbot will trigger a reimport of all 
the channel's pins. You can also trigger this manually with the `/import` command.

Don't forget that Pinbot needs [permission](#permissions) to see and post in these channels, otherwise it won't be able to do its job.

⚠️ Note that this bot is currently in [_beta_](https://github.com/elliotwms/pinbot/milestone/2). There may be bugs, please [report them](https://github.com/elliotwms/pinbot/issues/new?labels=bug&template=bug_report.md) ⚠️

### How does it _really_ work? Like, under the hood

Pinbot is deployed as an AWS Lambda function

#### Permissions

Pinbot is designed to be run with as few permissions as possible, however as part of its core functionality it needs to 
be able to read the contents of messages in your server. If you're not cool with this then you're welcome to audit the
code yourself, or [host and run your own Pinbot](#run).

Pinbot requires the following permissions to function in any channels you intend to use it:
* Read messages (`VIEW_CHANNEL`)
* Send messages (`SEND_MESSAGES`)
* Add reactions (`ADD_REACTIONS`)

## Development

### Configuration

| Variable             | Description                                                                                          | Required |
|----------------------|------------------------------------------------------------------------------------------------------|----------|
| `DISCORD_TOKEN`      | Bot token                                                                                            | `true`   |
| `DISCORD_PUBLIC_KEY` | Bot public key                                                                                       | `true`   |
| `LOG_LEVEL`          | [Log level](https://github.com/sirupsen/logrus#level-logging). `trace` enables discord-go debug logs | `false`  |

## Testing

`/tests` contains a suite of integration tests which run against [fakediscord](https://github.com/elliotwms/fakediscord) in a test guild. Simply run `docker-compose up` from the root of the repo and execute the tests.
