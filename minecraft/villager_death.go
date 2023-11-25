package minecraft

import (
	"regexp"
	"strconv"
)

type villagerDeathRegexpHandler struct {
	Base *regexp.Regexp
	UUID *regexp.Regexp
	L    *regexp.Regexp
	X    *regexp.Regexp
	Y    *regexp.Regexp
	Z    *regexp.Regexp
	V    *regexp.Regexp
}

var villagerDeathRegexp = villagerDeathRegexpHandler{
	Base: regexp.MustCompile(`Villager\sEntityVillager\['(.+)'\/(\d*)(.*)\]\sdied,\smessage:\s'(.+)'`),
	UUID: regexp.MustCompile(`uuid='([\S]+)'`),
	L:    regexp.MustCompile(`l='(\S+)'`),
	X:    regexp.MustCompile(`x=(-?[\d\.]+)`),
	Y:    regexp.MustCompile(`y=(-?[\d\.]+)`),
	Z:    regexp.MustCompile(`z=(-?[\d\.]+)`),
	V:    regexp.MustCompile(`v=(true|false)`),
}

func (h villagerDeathRegexpHandler) Match(message string) bool {
	return h.Base.MatchString(message)
}

func (h villagerDeathRegexpHandler) Parse(message string) MessageContentVillagerDeath {
	var result MessageContentVillagerDeath

	match := h.Base.FindAllStringSubmatch(message, 1)
	if len(match) == 0 {
		return result
	}

	result.Job = match[0][1]
	result.ID, _ = strconv.Atoi(match[0][2])
	result.DiedMessage = match[0][4]

	match = h.UUID.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.UUID = match[0][1]
	}

	match = h.L.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.L = match[0][1]
	}

	match = h.X.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.X, _ = strconv.ParseFloat(match[0][1], 64)
	}

	match = h.Y.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.Y, _ = strconv.ParseFloat(match[0][1], 64)
	}

	match = h.Z.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.Z, _ = strconv.ParseFloat(match[0][1], 64)
	}

	match = h.V.FindAllStringSubmatch(message, 1)
	if len(match) > 0 {
		result.V = match[0][1] == "true"
	}

	return result
}
