package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cufee/botto-lfg/config"
	"github.com/cufee/botto-lfg/utils"
)

var token string

// Channels
var channelNames map[string][]int = make(map[string][]int)

// Enebled cats
var enabledCats []string = []string{"809954293422751794"}

func init() {
	var err error
	token, err = utils.LoadToken("config/token.dat")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Event handlers

	// Add Voice channel handler
	dg.AddHandler(channelJoinedHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Botto LFG is now running. Press CTRL-C to exit.")
	// Wait for the user to cancel the process
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		dg.Close()
	}()
}

func channelJoinedHandler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// Get guild from State cache
	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		log.Print("guild not found in state")
		return
	}

	// Map to store member count per channel
	var validChannels map[string][]*discordgo.Channel = make(map[string][]*discordgo.Channel)

	// Get a list of channels
	for _, channel := range guild.Channels {
		// Check type and cetegory ID
		if channel.Type != discordgo.ChannelTypeGuildVoice || !utils.StringInSlice(channel.ParentID, enabledCats) {
			continue
		}

		// Save channel data
		validChannels[channel.ParentID] = append(validChannels[channel.ParentID], channel)
	}

	// Find empty channels
	for _, state := range guild.VoiceStates {
		// Get channel
		channel, err := s.Channel(e.ChannelID)
		if err != nil {
			continue
		}

		// Check if channel is in a right cat
		if !utils.StringInSlice(channel.ParentID, enabledCats) {
			continue
		}

		// Get category channels
		catChannels, ok := validChannels[channel.ParentID]
		if !ok {
			continue
		}

		// Pop channel from list
		for i, c := range catChannels {
			if c.ID == state.ChannelID {
				validChannels[channel.ParentID] = append(catChannels[:i], catChannels[i+1:]...)
			}
		}
	}

	// Check if channels need to be added or deleted
	for cat, emptyChannels := range validChannels {
		// Sort channels by position
		emptyChannels = utils.QuickSort(emptyChannels)

		// Delete extra channels
		if len(emptyChannels) > config.FreeChannelsBuffer {
			for i := 0; i < (len(emptyChannels) - config.FreeChannelsBuffer); i++ {
				_, err := s.ChannelDelete(emptyChannels[i].ID)
				if err != nil {
					log.Printf("failed to delete a channel: %v", err)
				}
			}
			continue
		}

		// Add free channels
		if len(emptyChannels) < config.FreeChannelsBuffer {
			for i := 0; i < (config.FreeChannelsBuffer - len(emptyChannels)); i++ {
				var chanData discordgo.GuildChannelCreateData
				chanData.UserLimit = config.FreeChannelsUserLimit
				chanData.Type = discordgo.ChannelTypeGuildVoice
				chanData.Name = fmt.Sprintf("LFG #%v", 0)
				chanData.ParentID = cat
				chanData.Position = 0
				_, err := s.GuildChannelCreateComplex(e.GuildID, chanData)
				if err != nil {
					log.Printf("failed to create a channel: %v", err)
				}
			}
			continue
		}
	}
}
