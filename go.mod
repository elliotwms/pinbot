module github.com/elliotwms/pinbot

go 1.23

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/bwmarrin/discordgo v0.28.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/elliotwms/fakediscord v0.12.18
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	github.com/winebarrel/secretlamb v0.3.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bwmarrin/discordgo => ../../bwmarrin/discordgo
