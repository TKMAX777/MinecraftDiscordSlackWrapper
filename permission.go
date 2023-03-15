package main

import "regexp"

type Permission string

func (p Permission) Verify(str string) (bool, error) {
	return regexp.MatchString("^"+string(p), str)
}

type Permissions []Permission

func (ps Permissions) Verify(str string) (bool, error) {
	for _, p := range ps {
		ok, err := p.Verify(str)
		if err != nil {
			return false, err
		}
		if ok {
			return ok, nil
		}
	}
	return false, nil
}
