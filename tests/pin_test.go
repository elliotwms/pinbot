package tests

import (
	"testing"
)

func TestPin(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_successfully_acknowledge_the_pin()
}

func TestPinGeneralPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		the_message_is_posted()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().
		the_bot_should_successfully_acknowledge_the_pin()
}

func TestPinSpecificPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		a_channel_named("test-pins").and().
		the_message_is_posted()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_successfully_acknowledge_the_pin()
}

func TestPinAlreadyPinned(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted().and().
		the_message_is_already_marked_as_pinned()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		the_bot_should_respond_with_message_containing("🔄 Message already pinned")
}

func TestPinWithImage(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_message().and().
		an_image_attachment().and().
		the_message_is_posted().and().
		the_message_has_n_attachments(1)

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_successfully_acknowledge_the_pin().and().
		the_pin_message_should_have_n_embeds(1).and().
		the_pin_message_should_have_an_image_embed()
}

func TestPinWithMultipleImage(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_message().and().
		an_image_attachment().and().
		another_image_attachment().and().
		the_message_is_posted()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_successfully_acknowledge_the_pin().and().
		the_pin_message_should_have_n_embeds(2).and().
		the_pin_message_should_have_n_embeds_with_image_url(2)
}

func TestPinWithFile(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_message().and().
		a_file_attachment().and().
		the_message_is_posted()

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_successfully_acknowledge_the_pin().and().
		the_pin_message_should_have_n_embeds(1).and().
		the_pin_message_should_have_n_embeds_with_image_url(0)
}

func TestPinPersistsEmbeds(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_message().and().
		the_message_has_a_link().and(). // posting a message with a link will create an embed on the server-side
		the_message_is_posted().and().
		the_message_has_n_embeds(1) // account for delay in link embed arriving (via MESSAGE_UPDATE)

	when.
		the_pin_command_is_sent_for_the_message()

	then.
		the_bot_should_successfully_acknowledge_the_pin().and().
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_pin_message_should_have_n_embeds(2) // the pin embed + link
}
