package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// DB - Database type
var DB map[string][]string = make(map[string][]string)

// DBPath - Path to database
const DBPath string = "database/database.json"

func readStorage() error {
	// Open file
	data, err := os.Open(DBPath)
	if err != nil {
		data, err = os.Create(DBPath)
		if err != nil {
			return err
		}
	}
	defer data.Close()

	// Unmarshal
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &DB)
	return nil
}

func writeStorage() error {
	// Open file
	data, err := os.Open(DBPath)
	if err != nil {
		data, err = os.Create(DBPath)
		if err != nil {
			return err
		}
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
