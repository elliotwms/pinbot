# Pinbot

[Install 📌](https://discord.com/oauth2/authorize?client_id=921554139740254209&permissions=3136&redirect_uri=https%3A%2F%2Fgithub.com%2Felliotwms%2Fpinbot&scope=applications.commands%20bot)

[Join the Pinbot Discord Server](https://discord.gg/a3u2PZ6V28)

Whenever you react to a message with a pushpin 📌 emoji, Pinbot posts the message to a channel.

![Example of a Pinbot message](https://user-images.githubusercontent.com/4396779/147515477-850ab41a-6a89-4746-9f65-e27c259f7602.png)

Pinbot is designed as an extension to Discord's channel pins system. Use Pinbot to:
* Bypass Discord's 50-pin limit and create a historic stream of all your pins
* Collect all your server's pins into one place (with optional overrides)
* Give your server's pins a more permanent home

### Why does this exist?

Discord guilds use pins for a lot more than just highlighting important information. In fact, many guilds use the pin system as a form of memorialising a good joke, a savage putdown, or other memorable moments. As a result, the 50 pin per channel limit means that in order to keep something, you will eventually have to get rid of something else.

### How does it work?

Pinbot uses the channel name to decide where it will post. In order of priority it will pin in:
1. `#{channel}-pins`, where `channel` is the name of the channel the message was pinned in
2. `#pins`, a general pins channel
3. `#{channel}`, the channel the pin was posted in, so that if you don't want a separate pins channel you can instead 
search for pins by @pinbot in the channel

Whenever Pinbot pins a message, or whenever you update the actual channel pins, Pinbot will trigger a reimport of all 
the channel's pins. You can also trigger this manually with the `/import` command.

Don't forget that pinbot needs [permission](#permissions) to see and post in these channels, otherwise it won't be able to do its job.

⚠️ Note that this bot is currently in [_beta_](https://github.com/elliotwms/pinbot/milestone/2). There may be bugs, please [report them](https://github.com/elliotwms/pinbot/issues/new?labels=bug&template=bug_report.md) ⚠️

#### Emojis

Pinbot will react with the following emojis to provide feedback:

| Emoji | Meaning                                                                                                            |
|-------|--------------------------------------------------------------------------------------------------------------------|
| 👀    | Pinbot has seen your message and is currently processing it                                                        |
| ✅     | Pinbot has successfully pinned your message                                                                        |
| 💩    | Pinbot could not perform an action for an unspecified reason                                                       |
| 🔄    | Pinbot could not pin their own message. Pinbot hates recursion                                                     |
| 🚫    | Pinbot could not pin this message as it was in an excluded channel (only really applicable to self-hosted Pinbots) |

#### Permissions

Pinbot is designed to be run with as few permissions as possible, however as part of its core functionality it needs to 
be able to read the contents of messages in your server. If you're not cool with this then you're welcome to audit the
code yourself, or [host and run your own Pinbot](#run).

Pinbot requires the following permissions to function in any channels you intend to use it:
* Read messages (`VIEW_CHANNEL`)
* Send messages (`SEND_MESSAGES`)
* Add reactions (`ADD_REACTIONS`)

## Run

Pinbot is designed to be run as the managed application above, but if you prefer (or if you don't trust a bot with 
permission to read and relay your messages) you can run your own. You will need to [create a new bot](https://discord.com/developers/applications),
obtain the token and application ID, and install the bot to your server (Pinbot will output a link to install the bot to
your servers in the "Starting Pinbot" message when it's run).

Part of the build pipeline includes building a Docker image which is [pushed to ghcr](https://github.com/elliotwms/pinbot/pkgs/container/pinbot).

```shell
export TOKEN {bot_token}
export APPLICATION_ID {bot_application_id}
docker run -e TOKEN -e APPLICATION_ID ghcr.io/elliotwms/pinbot:{version}
```

### Configuration

| Variable            | Description                                                                                          | Required |
|---------------------|------------------------------------------------------------------------------------------------------|----------|
| `TOKEN`             | Bot token ID                                                                                         | `true`   |
| `APPLICATION_ID`    | Bot application ID                                                                                   | `true`   |
| `TEST_GUILD_ID`     | When specified, the bot should only respond to pins in this test guild                               | `false`  |
| `HEALTH_CHECK_ADDR` | Address to serve the `/v1/health/` endpoint on (e.g. `:8080`)                                        | `false`  |
| `EXCLUDED_CHANNELS` | Comma-separated list of excluded channel IDs                                                         | `false`  |
| `LOG_LEVEL`         | [Log level](https://github.com/sirupsen/logrus#level-logging). `trace` enables discord-go debug logs | `false`  |

## Testing

`/tests` contains a suite of integration tests which run against the real Discord API in a test guild. It will create
and destroy a guild for the test run (unless you specify a `TEST_GUILD_ID`, in which case it will use an existing one,
[as explained below](#debugging)).

In order to run these test yourself you will need to:

* [Create a new bot](https://discord.com/developers/applications), obtaining the bot token
* Run the tests with the `TOKEN` and `APPLICATION_ID` environment variables 

Set the `TEST_DEBUG` environment variable in order to enable debug logging in tests.

### Debugging

If you would rather view the bot activity, then it's also possible to use an existing guild instead of creating and
destroying one for each test run.

* Invite your bot to your guild, giving it the usual bot permissions as well additional permissions required to run the 
tests. A link will be output at the start of the tests to install the bot which requests the base permissions along with 
the following:
    * Manage Channels (to create channels during tests)
    * Read Messages (to assert on message creation)
    * Manage messages (to pin messages)
* Set the `TEST_GUILD_ID` environment variable to the test guild's ID when running the tests [as above](#testing)
