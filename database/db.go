package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// DB - Database type
var DB map[string][]string = make(map[string][]string)

// DBPath - Path to database
const DBPath string = "database.json"

func readStorage() error {
	// Open file
	data, err := os.OpenFile(DBPath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer data.Close()

	// Unmarshal
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &DB)
	return nil
}

func writeStorage() error {
	// Open file
	data, err := os.OpenFile(DBPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer data.Close()

	// Marshal
	file, err := json.MarshalIndent(DB, "", " ")
	if err != nil {
		return err
	}

	// Write to file
	err = ioutil.WriteFile(DBPath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

// EnableGuildCategory -
func EnableGuildCategory(guildID string, catID string) error {
	// Read localDB
	err := readStorage()
	if err != nil {
		return err
	}

	// Check if category already added
	for _, c := range DB[guildID] {
		if c == catID {
			return nil
		}
	}

	// Add category
	DB[guildID] = append(DB[guildID], catID)

	// Write localDB
	err = writeStorage()
	if err != nil {
		return err
	}
	return nil
}

// DisableGuildCategory -
func DisableGuildCategory(guildID string, catID string) error {
	// Read localDB
	err := readStorage()
	if err != nil {
		return err
	}

	// Remove from slice
	for i, cid := range DB[guildID] {
		if cid == catID {
			DB[guildID] = append(DB[guildID][:i], DB[guildID][i+1:]...)
		}
	}

	// Write localDB
	err = writeStorage()
	if err != nil {
		return err
	}
	return nil
}

// GetGuildCategories -
func GetGuildCategories(guildID string) (cats []string) {
	readStorage()
	return DB[guildID]
}
