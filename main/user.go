package main

import (
	"encoding/json"
	"io/ioutil"
)

// Users is a format of a user list
type Users []User

// User is a format of a user profile
type User struct {
	DiscordID      string
	Name           string
	PermissionCode uint64
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

func (u Users) findUserFromDiscordID(id string) (User, bool) {
	for _, user := range u {
		if user.DiscordID == id {
			return user, true
		}
	}
	return User{}, false
}
