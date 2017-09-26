package core

import (
	"log"

	"github.com/larsfox/newochem-bot/db"
	"github.com/larsfox/newochem-bot/tg"
	"github.com/larsfox/newochem-bot/vk"
)

// Manager works with other packages
type Manager interface {
	Listen()
}

type manager struct {
	dbClient db.Client
	tgClient tg.Client
	tgUsers  []string
	vkClient vk.Client
}

// NewManager returns a new manager
func NewManager(dbClient db.Client, vkClient vk.Client, tgClient tg.Client, tgUsers []string) Manager {
	return &manager{
		dbClient: dbClient,
		tgClient: tgClient,
		tgUsers:  tgUsers,
		vkClient: vkClient,
	}
}

// Listen listens the incoming messages
func (m *manager) Listen() {
	log.Println("Listening...")

	/*
		TESTING AREA STARTS
	*/

	/*
		TESTING AREA ENDS
	*/

	msgChan, _ := m.tgClient.GetMessagesChan()
	for msg := range msgChan {
		if msg != nil {
			go m.handleMsg(msg)
		}
	}
}

func (m *manager) SendError(chatID int, err error) {
	log.Println(err)
	m.tgClient.SendMessage(chatID, errorString, "Markdown", nil)
	// bugsnag.Notify(err)
}
