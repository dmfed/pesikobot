package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

var (
	commandRegexp = regexp.MustCompile(`/(\w+) *`)
)

//go:embed hello.txt
var helloMsg string

type botOwner struct {
	id       int
	idstr    string
	username string
}

func (b botOwner) Recipient() string {
	return b.idstr
}

type pesikobot struct {
	*telebot.Bot
	botOwner
}

func (pes *pesikobot) all(m *telebot.Message) {
	if !pes.authorized(m) {
		pes.Send(m.Chat, "You are not authorized to issue commands. Sorry...")
		return
	}
	command := ""
	if commandRegexp.MatchString(m.Text) {
		command = commandRegexp.FindStringSubmatch(m.Text)[1]
		loc := commandRegexp.FindStringIndex(m.Text)
		m.Text = m.Text[loc[1]:]
	}
	switch command {
	case "pic":
		pes.picture(m)
	case "cmd":
		pes.cmd(m)
	case "whoami":
		pes.whoami(m)
	case "help":
		pes.help(m)
	}
}

func (pes *pesikobot) authorized(m *telebot.Message) bool {
	return m.Sender.Username == pes.username || m.Sender.ID == pes.id
}

func (pes *pesikobot) help(m *telebot.Message) {
	pes.send(m, helloMsg)
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

func (pes *pesikobot) whoami(m *telebot.Message) {
	var b bytes.Buffer
	usr := m.Sender
	b.WriteString(fmt.Sprintf("Your user id is: %d\n", usr.ID))
	b.WriteString(fmt.Sprintf("Username: %v\n", usr.Username))
	b.WriteString(fmt.Sprintf("First name: %v\nLast name: %v\n", usr.FirstName, usr.LastName))
	pes.Send(m.Chat, b.String())
}

func (pes *pesikobot) cmd(m *telebot.Message) {
	commslice := strings.Split(m.Text, " ")
	var cmd *exec.Cmd
	if len(commslice) > 1 {
		cmd = exec.Command(commslice[0], commslice[1:]...)
	} else {
		cmd = exec.Command(commslice[0])
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		pes.send(m, fmt.Sprintf("cmd.CombinedOutput() returned: %v", err.Error()))
	}
	pes.send(m, string(output))
}

func (pes *pesikobot) send(m *telebot.Message, msg string) {
	messages := paginate(string(msg))
	for _, message := range messages {
		pes.Send(m.Chat, message)
	}
}

func paginate(input string) []string {
	out := []string{}
	var buf bytes.Buffer
	for _, s := range strings.Split(input, "\n") {
		if buf.Len()+len(s) > 4096 {
			out = append(out, buf.String())
			buf.Reset()
		}
		buf.WriteString(s + "\n")
	}
	out = append(out, buf.String())
	return out
}
