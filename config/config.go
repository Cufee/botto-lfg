package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Data - Config data struct
type Data struct {
	Prefix       string `json:"prefix"`
	ChannelsBuff int    `json:"channels_buffer"`
	UserLimit    int    `json:"channel_user_limit"`
	EventSpacing int    `json:"event_spacing"`
	Token        string `json:"token"`
}

// Defaults - Fill config with deaults
func (c *Data) Defaults() {
	c.Prefix = "b-"
	c.UserLimit = 0
	c.ChannelsBuff = 3
	c.EventSpacing = 5
	c.Token = "place_your_token_here"
}

// Read - Read local config file
func Read() (currentConfig Data) {
	// Open file
	data, err := os.OpenFile("config.json", os.O_CREATE, 0644)
	if err != nil {
		data, err = os.Create("config.json")
	}
	defer data.Close()

	// Unmarshal
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &currentConfig)

	// Check if config is empty
	if (currentConfig) == (Data{}) {
		// Marshal
		currentConfig.Defaults()
		file, err := json.MarshalIndent(&currentConfig, "", " ")
		if err != nil {
			panic(err)
		}

		// Write
		_, err = data.Write(file)
		if err != nil {
			panic(err)
		}
	}

	return currentConfig
}
