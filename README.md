# Pinbot

> [!TIP]
> Looking for a more lightweight Pinbot which runs as a Lambda function? check out [elliotwms/pinbot-lambda](https://github.com/elliotwms/pinbot-lambda)

[Install 📌](https://discord.com/discovery/applications/921554139740254209)

[Need help? Join the Pinbot Discord Server](https://discord.gg/a3u2PZ6V28)

Pinbot provides a single message command, Pin, which posts a copy of your message to your server's pins channel.

![Example of a Pinbot message](https://user-images.githubusercontent.com/4396779/147515477-850ab41a-6a89-4746-9f65-e27c259f7602.png)

Pinbot is designed as an extension to Discord's channel pins system. Use Pinbot to:
* Bypass Discord's 50-pin limit and create a historic stream of all your pins
* Collect all your server's pins into one place (with optional overrides)
* Give your server's pins a more permanent home

### Why does this exist?

Discord guilds use pins for a lot more than just highlighting important information. In fact, many guilds use the pin system as a form of memorialising a good joke, a savage putdown, or other memorable moments. As a result, the 50 pin per channel limit means that in order to keep something, you will eventually have to get rid of something else.

### How does it work?

Pinbot uses the channel name of the pinned message to decide where it will post. In order of priority it will pin in:
1. `#{channel}-pins`, where `channel` is the name of the channel the message was pinned in
2. `#pins`, a general pins channel
3. `#{channel}`, the channel the pin was posted in, so that if you don't want a separate pins channel you can instead search for pins by @pinbot in the channel

Don't forget that pinbot needs [permission](#permissions) to see and post in these channels, otherwise it won't be able to do its job.

⚠️ Note that this bot is currently in [_beta_](https://github.com/elliotwms/pinbot/milestone/2). There may be bugs, please [report them](https://github.com/elliotwms/pinbot/issues/new?labels=bug&template=bug_report.md) ⚠️

#### Permissions

Pinbot is designed to be run with as few permissions as possible, however as part of its core functionality it needs to be able to read the contents of messages in your server. If you're not cool with this then you're welcome to audit the code yourself, or [host and run your own Pinbot](#run).

Pinbot requires the following permissions to function in any channels you intend to use it:
* Read messages (`VIEW_CHANNEL`)
* Send messages (`SEND_MESSAGES`)
* Add reactions (`ADD_REACTIONS`)

## Run

Pinbot is designed to be run as the managed application above, but if you prefer (or if you don't trust a bot with permission to read and relay your messages) you can run your own. You will need to [create a new bot](https://discord.com/developers/applications), obtain the token and application ID, and install the bot to your own server.

Part of the build pipeline includes building a Docker image which is [pushed to ghcr](https://github.com/elliotwms/pinbot/pkgs/container/pinbot).

```shell
export TOKEN {bot_token}
export APPLICATION_ID {bot_application_id}
docker run -e TOKEN -e APPLICATION_ID ghcr.io/elliotwms/pinbot:{version}
```

### Configuration

| Variable            | Description                                                                                                                   | Required |
|---------------------|-------------------------------------------------------------------------------------------------------------------------------|----------|
| `TOKEN`             | Bot token ID                                                                                                                  | `true`   |
| `APPLICATION_ID`    | Bot application ID                                                                                                            | `true`   |
| `GUILD_ID`          | When specified, the bot should only migrate commands within this guild. Useful for testing or running your own solo-guild bot | `false`  |
| `HEALTH_CHECK_ADDR` | Address to serve the `/v1/health/` endpoint on (e.g. `:8080`)                                                                 | `false`  |
| `LOG_LEVEL`         | [Log level](https://pkg.go.dev/log/slog#hdr-Levels). `debug` enables discord-go debug logs                                    | `false`  |

## Testing

`/tests` contains a suite of integration tests which run against fakediscord
