package config

// BotPrefix - valid prefix for the bot
var BotPrefix string = "b-"

// FreeChannelsBuffer - how many free channels to always keep available
const FreeChannelsBuffer int = 3

// FreeChannelsUserLimit - how many free slots to make free channels with
const FreeChannelsUserLimit int = 0 // If set to 0, largest userLimit in the category will be used

// UpdateDelataySec - how long to wait after a user leaves a channel to delete/create empty channels, adding a delay will deduce spam
const UpdateDelataySec int = 5
