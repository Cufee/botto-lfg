package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cufee/botto-lfg/config"
	"github.com/cufee/botto-lfg/database"
	"github.com/cufee/botto-lfg/utils"
)

var token string
var eventChan = make(chan int, 1)

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
	dg.AddHandler(voiceEvents)

	// Enable/Disable bot for a category
	dg.AddHandler(addCatCommand)
	dg.AddHandler(removeCatCommand)

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

// addCatCommand - Add a category for the bot to watch
func addCatCommand(s *discordgo.Session, e *discordgo.MessageCreate) {
	var command string = (config.BotPrefix + "watchcat")
	// Check for prefix
	if !strings.HasPrefix(e.Content, command) {
		return
	}

	// Get category id
	catID := strings.TrimSpace(strings.ReplaceAll(e.Content, command, ""))
	channel, err := s.State.GuildChannel(e.GuildID, catID)
	if err != nil {
		s.ChannelMessageSend(e.ChannelID, "Failed to find a channel with this ID, does the bot have proper perms?")
		return
	}
	if channel.Type != discordgo.ChannelTypeGuildCategory {
		s.ChannelMessageSend(e.ChannelID, "This channel is not a category, please copy and ID from the category you want me to watch.")
		return
	}

	// Write to DB
	if err = database.EnableGuildCategory(e.GuildID, catID); err != nil {
		s.ChannelMessageSend(e.ChannelID, "Failed to add this category to my database.")
		return
	}

	s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("I am now watching %v.", channel.Name))
}

// removeCatCommand - Add a category for the bot to watch
func removeCatCommand(s *discordgo.Session, e *discordgo.MessageCreate) {
	var command string = (config.BotPrefix + "lookaway")
	// Check for prefix
	if !strings.HasPrefix(e.Content, command) {
		return
	}

	// Get category id
	catID := strings.TrimSpace(strings.ReplaceAll(e.Content, command, ""))
	channel, err := s.State.GuildChannel(e.GuildID, catID)
	if err != nil {
		s.ChannelMessageSend(e.ChannelID, "Failed to find a channel with this ID, does the bot have proper perms?")
		return
	}
	if channel.Type != discordgo.ChannelTypeGuildCategory {
		s.ChannelMessageSend(e.ChannelID, "This channel is not a category, please copy and ID from the category you want me to watch.")
		return
	}

	// Write to DB
	if err = database.DisableGuildCategory(e.GuildID, catID); err != nil {
		s.ChannelMessageSend(e.ChannelID, "Failed to remove this category from my database.")
		return
	}

	s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("I am no longer watching %v.", channel.Name))
}

// voiceEvents - Handler for voice stats updates
func voiceEvents(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// Get guild from State cache
	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		log.Print("guild not found in state")
		return
	}

	// Check for pending events
	select {
	case eventChan <- 1: // Put 1 in the channel unless it is full
	default:
		// Event pending, return
		return
	}
	defer func() { <-eventChan }()

	// Sleep to avoid spam
	if config.UpdateDelataySec > 0 {
		time.Sleep(time.Second * time.Duration(config.UpdateDelataySec))
	}

	// Map to store member count per channel
	var validChannels map[string][]*discordgo.Channel = make(map[string][]*discordgo.Channel)

	// Enebled cats
	var enabledCats []string = database.GetGuildCategories(e.GuildID)

	// Valid channel IDs
	var validNames map[int]bool = make(map[int]bool)

	// Channel name template
	var nameTemplate string

	// Get a list of channels
	for _, channel := range guild.Channels {
		// Check type and cetegory ID
		if channel.Type != discordgo.ChannelTypeGuildVoice || !utils.StringInSlice(channel.ParentID, enabledCats) {
			continue
		}

		// Check if channel is added already
		var skip bool
		for _, c := range validChannels[channel.ParentID] {
			if c.Name == channel.Name {
				skip = true
				break
			}
		}

		if !skip {
			// Set ID as taken
			var sep string = "#"
			nameSlice := strings.Split(channel.Name, sep)
			if len(nameSlice) == 1 {
				sep = " "
				nameSlice = strings.Split(channel.Name, sep)
			}
			channelNum, _ := strconv.Atoi(nameSlice[len(nameSlice)-1])
			validNames[channelNum] = true

			// Set name template
			if nameTemplate == "" {
				nameTemplate = nameSlice[0] + sep
			}

			// Save channel data
			validChannels[channel.ParentID] = append(validChannels[channel.ParentID], channel)
		}
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
			for i := len(emptyChannels) - 1; i >= (config.FreeChannelsBuffer); i-- {
				// Delete channel
				_, err := s.ChannelDelete(emptyChannels[i].ID)
				if err != nil {
					log.Printf("failed to delete a channel: %v", err)
					continue
				}
			}
			continue
		}

		// Add free channels
		if len(emptyChannels) < config.FreeChannelsBuffer {
			for i := 0; i < (config.FreeChannelsBuffer - len(emptyChannels)); i++ {
				// Find next available id
				var channelID int
				for i := 1; true; i++ {
					if !validNames[i] {
						channelID = i
						validNames[i] = true
						break
					}
				}

				// Create a channel
				var chanData discordgo.GuildChannelCreateData
				chanData.UserLimit = config.FreeChannelsUserLimit
				chanData.Type = discordgo.ChannelTypeGuildVoice
				chanData.Name = fmt.Sprintf("%v%v", nameTemplate, channelID)
				chanData.ParentID = cat
				_, err := s.GuildChannelCreateComplex(e.GuildID, chanData)
				if err != nil {
					log.Printf("failed to create a channel: %v", err)
				}
				validNames[channelID] = true
			}
			continue
		}
	}
}
