package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type UsersSettings struct {
	Users map[string]User `yaml:"Users"`

	// map[GroupName]PermissionRegexp
	Permissions map[string]Permissions `yaml:"Permissions"`
}

// GetUsersSettings reads user settings
func GetUsersSettings() (*UsersSettings, error) {
	var yamlRootPath = filepath.Join("settings", "users")

	dir, err := os.ReadDir(yamlRootPath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDir")
	}

	var yamlBinary = []byte{}
	for _, f := range dir {
		if f.IsDir() || !(strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
			continue
		}

		var yamlFilePath = filepath.Join(yamlRootPath, f.Name())
		b, err := os.ReadFile(yamlFilePath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("ReadFile: %s", yamlFilePath))
		}

		yamlBinary = append(yamlBinary, b...)
		yamlBinary = append(yamlBinary, '\n')
	}

	var us UsersSettings
	err = yaml.Unmarshal(yamlBinary, &us)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return &us, nil
}

type ServiceType int

const (
	ServiceTypeSlack ServiceType = iota
	ServiceTypeDiscord
)

func (us UsersSettings) GetUser(userid string, serviceType ServiceType) *User {
	for name, u := range us.Users {
		switch serviceType {
		case ServiceTypeDiscord:
			if u.Discord == userid {
				u.Groups = append(u.Groups, "everyone")
				u.Name = name
				return &u
			}
		case ServiceTypeSlack:
			if u.Slack == userid {
				u.Groups = append(u.Groups, "everyone")
				u.Name = name
				return &u
			}
		default:
			return nil
		}
	}
	return &User{
		Groups:    []string{"everyone"},
		IsDefault: true,
	}
}

func (us UsersSettings) GetPermissions(groups []string) Permissions {
	var permissions = Permissions{}
	for g, ps := range us.Permissions {
		for _, ug := range groups {
			if ug == g {
				permissions = append(permissions, ps...)
			}
		}
	}
	return permissions
}
