package main

import "strings"

type params struct {
	params  map[string]string
	remote  string
	org     string
	user    string
	repo    string
	goos    string
	goarch  string
	version string
	token   string
}

func newParams(p map[string]string, token string) *params {
	return &params{
		params:  p,
		remote:  p["remote"],
		org:     p["org"],
		user:    p["user"],
		repo:    p["repo"],
		goos:    p["goos"],
		goarch:  p["goarch"],
		version: p["version"],
		token:   parseToken(token),
	}
}

func (p params) owner() string {
	var owner string
	if p.org != "" {
		owner = p.org + "/"
	}
	return owner + p.user
}

func parseToken(token string) string {
	splited := strings.Split(token, " ")
	if len(splited) > 1 {
		return splited[1]
	}
	return token
}
