package tests

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/endpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PinStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	assert  *assert.Assertions

	handler func(_ context.Context, event *events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLResponse, error)
	res     *events.LambdaFunctionURLResponse

	sendMessage         *discordgo.MessageSend
	channel             *discordgo.Channel
	expectedPinsChannel *discordgo.Channel
	message             *discordgo.Message

	messages            []*discordgo.Message
	pinMessage          *discordgo.Message
	snowflake           *snowflake.Node
	interactionResponse *discordgo.InteractionResponse
}

func NewPinStage(t *testing.T) (*PinStage, *PinStage, *PinStage) {
	node, _ := snowflake.NewNode(0)
	s := &PinStage{
		t:         t,
		session:   session,
		require:   require.New(t),
		assert:    assert.New(t),
		handler:   endpoint.New([]byte{}, "").WithSession(session).WithApplicationCommand("Pin", commandhandlers.PinMessageCommandHandler).Handle,
		snowflake: node,
	}

	_, cancel := context.WithCancel(context.Background())

	t.Cleanup(cancel)

	return s, s, s
}

func (s *PinStage) and() *PinStage {
	return s
}

func (s *PinStage) a_channel_named(name string) *PinStage {
	c, err := s.session.GuildChannelCreate(testGuildID, name, discordgo.ChannelTypeGuildText)
	s.require.NoError(err)

	s.t.Cleanup(func() {
		_, err = s.session.ChannelDelete(c.ID)
		s.assert.NoError(err)
	})

	if s.channel == nil {
		// register the first created channel as the "default" channel for the stage
		s.channel = c
	}
	// register the last created channel as the expected pins channel
	s.expectedPinsChannel = c

	s.session.AddHandler(s.handleMessageFor(c.ID))

	return s
}

func (s *PinStage) a_message() *PinStage {
	s.sendMessage = &discordgo.MessageSend{
		Content: "Hello, World!",
	}

	return s
}

func (s *PinStage) the_message_is_posted() *PinStage {
	if s.sendMessage == nil {
		s.a_message()
	}

	m, err := s.session.ChannelMessageSendComplex(s.channel.ID, s.sendMessage)
	s.require.NoError(err)
	s.message = m

	return s
}

func (s *PinStage) the_pin_command_is_sent_for_the_message() *PinStage {
	bs, err := json.Marshal(&discordgo.InteractionCreate{
		&discordgo.Interaction{
			ID:    s.snowflake.Generate().String(),
			AppID: "1290742494824366183",
			Type:  discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				ID:          s.snowflake.Generate().String(), // todo command ID
				Name:        "Pin",
				CommandType: discordgo.MessageApplicationCommand,
				TargetID:    s.message.ID,
				Resolved: &discordgo.ApplicationCommandInteractionDataResolved{
					Messages: map[string]*discordgo.Message{
						s.message.ID: s.message,
					},
				},
			},
			GuildID:        testGuildID,
			ChannelID:      s.message.ChannelID,
			AppPermissions: 0,
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID: s.snowflake.Generate().String(),
				},
			},
			Context: discordgo.InteractionContextGuild,
			Version: 1,
		},
	})
	s.require.NoError(err)

	s.res, err = s.handler(context.Background(), &events.LambdaFunctionURLRequest{
		RequestContext:  events.LambdaFunctionURLRequestContext{},
		Body:            string(bs),
		IsBase64Encoded: false,
	})
	if err != nil {
		return nil
	}

	s.interactionResponse = &discordgo.InteractionResponse{}
	s.require.NoError(json.Unmarshal([]byte(s.res.Body), s.interactionResponse))

	return s
}

func (s *PinStage) a_pin_message_should_be_posted_in_the_last_channel() *PinStage {
	s.require.Eventually(func() bool {
		s.t.Logf("msgs: %d", len(s.messages))
		for _, m := range s.messages {
			if m.ChannelID != s.expectedPinsChannel.ID {
				continue
			}

			for _, embed := range m.Embeds {
				if embed.Title == "📌 Pinned" && strings.Contains(embed.Description, s.sendMessage.Content) {
					s.pinMessage = m
					return true
				}
			}
		}

		return false
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *PinStage) the_bot_should_add_the_emoji(emoji string) *PinStage {
	s.require.Eventually(func() bool {
		reactions, err := s.session.MessageReactions(s.channel.ID, s.message.ID, emoji, 0, "", "")
		if err != nil {
			return false
		}

		for _, r := range reactions {
			if r.ID == s.session.State.User.ID {
				return true
			}
		}

		return false
	}, 5*time.Second, 500*time.Millisecond)

	return s
}

func (s *PinStage) handleMessageFor(channelID string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID == channelID {
			s.messages = append(s.messages, m.Message)
		}
	}
}

func (s *PinStage) the_message_is_already_marked_as_pinned() {
	s.require.NoError(s.session.MessageReactionAdd(s.message.ChannelID, s.message.ID, "📌"))
}

func (s *PinStage) the_bot_should_respond_with_message_containing(m string) *PinStage {
	s.require.Contains(s.interactionResponse.Data.Content, m)

	return s
}

func (s *PinStage) an_attachment(filename, contentType string) *PinStage {
	f, err := os.Open("files/" + filename)
	s.require.NoError(err)
	s.sendMessage.Files = append(s.sendMessage.Files, &discordgo.File{
		Name:        filename,
		ContentType: contentType,
		Reader:      f,
	})

	return s
}

func (s *PinStage) an_image_attachment() *PinStage {
	return s.an_attachment("cheese.jpg", "image/jpeg")
}

func (s *PinStage) another_image_attachment() *PinStage {
	return s.an_image_attachment()
}

func (s *PinStage) a_file_attachment() *PinStage {
	return s.an_attachment("hello.txt", "text/plain")
}

func (s *PinStage) the_pin_message_should_have_an_image_embed() {
	s.the_pin_message_should_have_n_embeds_with_image_url(1)
}

func (s *PinStage) the_pin_message_should_have_n_embeds_with_image_url(n int) {
	found := 0
	for _, embed := range s.pinMessage.Embeds {
		if embed.Image != nil && embed.Image.URL != "" {
			found++
		}
	}

	s.require.Equal(n, found)
}

func (s *PinStage) the_pin_message_should_have_n_embeds(n int) *PinStage {
	s.require.Len(s.pinMessage.Embeds, n)

	return s
}

func (s *PinStage) the_import_is_cleaned_up() *PinStage {
	s.a_pin_message_should_be_posted_in_the_last_channel()

	s.require.NoError(s.session.ChannelMessageDelete(s.pinMessage.ChannelID, s.pinMessage.ID))
	s.messages = []*discordgo.Message{}

	s.require.NoError(s.session.MessageReactionsRemoveAll(s.message.ChannelID, s.message.ID))

	return s
}

func (s *PinStage) the_message_has_a_link() *PinStage {
	s.sendMessage.Content = s.sendMessage.Content + " https://github.com/elliotwms/pinbot"

	return s
}

func (s *PinStage) the_message_has_n_embeds(n int) {
	s.require.Eventually(func() bool {
		m, err := s.session.ChannelMessage(s.channel.ID, s.message.ID)
		if err != nil {
			return false
		}

		return len(m.Embeds) == n
	}, 5*time.Second, 500*time.Millisecond)
}

func (s *PinStage) the_message_has_n_attachments(n int) {
	s.require.Eventually(func() bool {
		m, err := s.session.ChannelMessage(s.channel.ID, s.message.ID)
		if err != nil {
			return false
		}

		return len(m.Attachments) == n
	}, 5*time.Second, 500*time.Millisecond)
}

func (s *PinStage) the_bot_should_successfully_acknowledge_the_pin() *PinStage {
	return s.
		the_bot_should_add_the_emoji("📌").and().
		the_bot_should_respond_with_message_containing("📌 Pinned")
}
