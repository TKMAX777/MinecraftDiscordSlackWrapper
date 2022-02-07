package minecraft

import (
	"log"
	"regexp"
	"strings"
)

type DeathMesasgeHandler []*regexp.Regexp

func NewDeathMesasgeHandler() (DeathMesasgeHandler, error) {
	var regexps = []*regexp.Regexp{}

	var lines = strings.Split(deathRegexps, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		reg, err := regexp.Compile(l)
		if err != nil {
			log.Printf("ParseDeathMessageErr: %s", err.Error())
			continue
		}
		regexps = append(regexps, reg)
	}

	return regexps, nil
}

func (d DeathMesasgeHandler) Match(message string) bool {
	for _, re := range d {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}
