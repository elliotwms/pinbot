package tests

import (
	"os"
	"testing"
)

func TestPin(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji()
}

func TestPinGeneralPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().
		the_bot_should_react_with_successful_emoji()
}

func TestPinSpecificPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		a_channel_named("test-pins").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji()
}

func TestPinAlreadyPinned(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted().and().
		the_message_is_already_marked_as_pinned()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		the_bot_should_log_the_message_as_already_pinned()
}

// TestPinSelfPinDisabled tests the 'correct' behaviour of Pinbot when pinning its own messages.
// Self-pin is enabled to allow testing via a single bot in a single server
func TestPinSelfPinDisabled(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		self_pin_is_disabled().and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		the_bot_should_add_the_emoji("🔄")
}

func TestPinClassicPinTriggersChannelImport(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted()

	when.
		the_message_is_pinned()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji()
}

func TestPinImportCommand(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted().and().
		the_message_is_pinned().and().
		the_bot_should_react_with_successful_emoji().and().
		the_import_is_cleaned_up()

	when.
		an_import_is_triggered()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji()
}

func TestPinImportCommandIgnoreAlreadyPinned(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted().and().
		the_message_is_pinned().and().
		the_bot_should_react_with_successful_emoji()

	when.
		an_import_is_triggered()

	then.
		the_bot_should_log_the_message_as_already_pinned()
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
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji().and().
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
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji().and().
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
		the_message_is_reacted_to_with("📌")

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_react_with_successful_emoji().and().
		the_pin_message_should_have_n_embeds(1).and().
		the_pin_message_should_have_n_embeds_with_image_url(0)
}

func TestPinInExcludedChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_channel_is_excluded().and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to_with("📌")

	then.
		the_bot_should_add_the_emoji("👀").and().
		the_bot_should_add_the_emoji("🚫")
}

func TestPinPersistsEmbeds(t *testing.T) {
	if os.Getenv("FAKEDISCORD") != "" {
		t.Skip("test incompatible with fakediscord: todo url should create server-side embed")
	}

	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_message().and().
		the_message_has_a_link().and(). // posting a message with a link will create an embed on the server-side
		the_message_is_posted().and().
		the_message_has_n_embeds(1) // account for delay in link embed arriving (via MESSAGE_UPDATE)

	when.
		the_message_is_reacted_to_with("📌")

	then.
		the_bot_should_react_with_successful_emoji().and().
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_pin_message_should_have_n_embeds(2) // the pin embed + link
}
