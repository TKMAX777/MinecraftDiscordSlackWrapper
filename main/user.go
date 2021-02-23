package main

import (
	"encoding/json"
	"io/ioutil"
)

// Users is a format of a user list
type Users []struct {
	DiscordID string
	Name      string
}

// ReadNameDict read name list
func ReadNameDict() (Users, error) {
	var userDict Users
	b, err := ioutil.ReadFile(NameDictPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &userDict)
	if err != nil {
		return nil, err
	}

	return userDict, nil
}

func (u Users) findNameFromDiscordID(id string) (string, bool) {
	for _, user := range u {
		if user.DiscordID == id {
			return user.Name, true
		}
	}
	return "", false
}
