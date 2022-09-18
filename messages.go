package goeinstein

import (
	"bufio"
	"bytes"
	"log"
	"strings"
)

func msg(key string) string {
	return messages.messages[key]
}

var messages = NewMessages()

type Messages struct {
	messages map[string]string
}

func NewMessages() *Messages {
	m := &Messages{}
	m.Load()
	return m
}

func (m *Messages) Load() {
	m.messages = make(map[string]string)

	bs := resources.GetRef("messages.txt")

	s := bufio.NewScanner(bytes.NewReader(bs))
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			log.Printf("Wrong line of messages: %q", line)
			continue
		}

		key = strings.Trim(key, " \"")
		value = strings.Trim(value, " \"")
		m.messages[key] = value
	}
}
