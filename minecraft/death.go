package minecraft

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type DeathMesasgeHandler []*regexp.Regexp

func NewDeathMesasgeHandler() (DeathMesasgeHandler, error) {
	b, err := ioutil.ReadFile(deathMessageTxt)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDeathMessageTxt")
	}

	var regexps = []*regexp.Regexp{}

	var lines = strings.Split(string(b), "\n")
	for _, l := range lines {
		reg, err := regexp.Compile(strings.TrimSpace(l))
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
