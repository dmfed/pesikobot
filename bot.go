package main

import (
	_ "embed"
	"os"
	"regexp"

	"gopkg.in/tucnak/telebot.v2"
)

var (
	commandRegexp = regexp.MustCompile(`/(\w+) *`)
)

//go:embed hello.txt
var helloMsg string

type pesikobot struct {
	*telebot.Bot
}

func (pes *pesikobot) all(m *telebot.Message) {
	command := ""
	if commandRegexp.MatchString(m.Text) {
		command = commandRegexp.FindStringSubmatch(m.Text)[1]
		loc := commandRegexp.FindStringIndex(m.Text)
		m.Text = m.Text[loc[1]:]
	}
	switch command {
	case "pic":
		pes.picture(m)
	case "hello":
		pes.hello(m)
	}
}

func (pes *pesikobot) hello(m *telebot.Message) {
	pes.Send(m.Chat, helloMsg)
}

func (pes *pesikobot) picture(m *telebot.Message) {
	filename, err := takePhoto()
	defer os.Remove(filename)
	if err != nil {
		pes.Send(m.Chat, err.Error())
		return
	}
	photo := &telebot.Photo{File: telebot.FromDisk(filename)}
	pes.Send(m.Chat, photo)
}
