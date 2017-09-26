package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/larsfox/newochem-bot/db"
	"github.com/larsfox/newochem-bot/tg"
)

const (
	stateNothing = iota
	stateGettingPosts
	stateCheckingPost
	stateTranslators
	stateEditors
	stateCategories
)

// handleMsg is the main function for handling messages
func (m *manager) handleMsg(msg *tg.Message) {
	if msg == nil {
		m.SendError(msg.Chat.ID, errors.New("Nil pointer on msg"))
		return
	}

	// User has no access
	if stringInArray(msg.Chat.Username, m.tgUsers) == -1 {
		return
	}

	// Getting current state
	var state *db.State
	var input *db.StateInput
	state, input, err := m.dbClient.GetState(msg.Chat.Username)
	if err != nil {
		switch err.Error() {
		case notFound: // if user has no state, create one
			state, err = m.dbClient.CreateState(msg.Chat.Username)
			if err != nil {
				m.SendError(msg.Chat.ID, err)
				return
			}

		default:
			m.SendError(msg.Chat.ID, err)
			return
		}
	}

	// If operation cancel
	if msg.Text == strCancel {
		state.State = 0
		input = nil
		m.tgClient.SendMessage(msg.Chat.ID, replyCancel, markdown, defaultKeyboard)
		m.dbClient.SetState(state, input) // TODO: copy-paste
		return
	}

	switch state.State {
	// Nothing
	case stateNothing:
		switch msg.Text {
		case strAddArticle, strAddArticleShort:
			// TODO: get article from the wall
			input.Article = &db.Article{}
			state.State = stateGettingPosts
			m.defaultWorkers(msg, state, translator)

		default:
			m.nothing(msg)
			return
		}

	// Getting wall posts
	case stateGettingPosts:

	// Checking if the gotten post is correct
	case stateCheckingPost:

	// Adding translators
	case stateTranslators:
		switch msg.Text {
		case strDone:
			if !checkJob(input, translator) {
				m.tgClient.SendMessage(msg.Chat.ID, replyNoChosen, markdown, nil)
				return
			}

			state.State = stateEditors
			m.defaultWorkers(msg, state, editor)
			m.tgClient.SendMessage(msg.Chat.ID, replyNoChosen, markdown, nil)

		default:
			m.addWorkers(msg, state, input, translator)
		}

	// Adding editors
	case stateEditors:
		switch msg.Text {
		case strDone:
			if !checkJob(input, editor) {
				m.tgClient.SendMessage(msg.Chat.ID, replyNoChosen, markdown, nil)
				return
			}

			state.State = stateCategories
			categories, err := m.dbClient.GetCategories()
			if err != nil {
				m.SendError(msg.Chat.ID, err)
				return
			}

			buttons := make([]string, len(categories))
			for i, c := range categories {
				buttons[i] = strAddPlus + c.Name
			}
			m.tgClient.SendMessage(msg.Chat.ID, replyCategoryJob, markdown, createAddKeyboards(buttons, false))

		default:
			m.addWorkers(msg, state, input, editor)
		}

	// Adding category
	case stateCategories:
		switch msg.Text {
		case strDone:
			state.State = stateNothing
			m.saveArticle(msg, input)

		default:
			m.addCategories(msg, state, input)
		}

	default:
		m.SendError(msg.Chat.ID, errors.New("Unknown state"))
	}

	m.dbClient.SetState(state, input)
}

func (m *manager) nothing(msg *tg.Message) {
	m.tgClient.SendMessage(msg.Chat.ID,
		fmt.Sprintf(replyDefault, msg.Chat.FirstName, msg.Chat.LastName),
		markdown, defaultKeyboard)
}

func (m *manager) cancel(msg *tg.Message, state *db.State, input *db.StateInput) {
	state.State = 0
	input = nil
	m.tgClient.SendMessage(msg.Chat.ID, replyCancel, markdown, defaultKeyboard)
}

// Sends starting message and keyboard for both translators and editors
func (m *manager) defaultWorkers(msg *tg.Message, state *db.State, jobKind string) {
	workers, err := m.dbClient.GetWorkers()
	if err != nil {
		m.SendError(msg.Chat.ID, err)
	}
	buttons := make([]string, len(workers))
	for i, w := range workers {
		buttons[i] = strAddPlus + w.ShortName
	}
	m.tgClient.SendMessage(msg.Chat.ID, jobMessages[jobKind], markdown, createAddKeyboards(buttons, false))
}

// Adds both translators and editors based on input
func (m *manager) addWorkers(msg *tg.Message, state *db.State, input *db.StateInput, jobKind string) {
	selected := make([]int, len(input.Jobs))
	for i, job := range input.Jobs {
		selected[i] = job.UserID
	}

	workers, err := m.dbClient.GetWorkers()
	if err != nil {
		m.SendError(msg.Chat.ID, err)
	}

	wids := make([]int, len(workers))
	names := make([]string, len(workers))
	buttons := make([]string, len(workers))
	for i, w := range workers {
		wids[i] = w.VKID
		names[i] = w.ShortName
		buttons[i] = strAddPlus + w.ShortName
	}

	var vkID int
	vkID, err = strconv.Atoi(msg.Text)
	if err != nil {
		split := strings.Split(msg.Text, " ")
		if len(split) == 1 {
			m.tgClient.SendMessage(msg.Chat.ID, replyNoUser, markdown, nil)
			return
		}

		index := stringInArray(strings.Join(split[1:], " "), names)
		if index == -1 {
			m.tgClient.SendMessage(msg.Chat.ID, replyNoUser, markdown, nil)
			return
		}
		vkID = workers[index].VKID
	}

	index := intInArray(vkID, selected)
	if index == -1 {
		input.Jobs = append(input.Jobs, &db.Job{
			// ArticleID: ,
			UserID: vkID,
			Kind:   jobKind,
		})
	} else {
		input.Jobs = append(input.Jobs[:index], input.Jobs[index+1:]...)
	}

	// All added users have a minus button and a mention in reply
	reply := jobReplies[jobKind]
	for _, job := range input.Jobs {
		if job.Kind != jobKind {
			continue
		}

		index := intInArray(job.UserID, wids)
		if index == -1 {
			reply += fmt.Sprintf("\n%d", job.UserID)
		} else {
			reply += "\n" + names[index]
			buttons[index] = strAddMinus + names[index]
		}
	}

	m.tgClient.SendMessage(msg.Chat.ID, reply,
		markdown, createAddKeyboards(buttons, reply != jobReplies[jobKind]))
}

// Adds categories based on input
func (m *manager) addCategories(msg *tg.Message, state *db.State, input *db.StateInput) {
	var category string
	categories, err := m.dbClient.GetCategories()
	if err != nil {
		m.SendError(msg.Chat.ID, err)
		return
	}

	if strings.HasPrefix(msg.Text, strAddPlus) || strings.HasPrefix(msg.Text, strAddMinus) {
		category = strings.Join(strings.Split(msg.Text, " ")[1:], " ")
	} else {
		category = msg.Text
	}

	names := make([]string, len(categories))
	buttons := make([]string, len(categories))
	for i, cat := range categories {
		names[i] = cat.Name
		buttons[i] = strAddPlus + cat.Name
	}

	index := stringInArray(category, names)
	if index == -1 {
		m.tgClient.SendMessage(msg.Chat.ID, replyNoCategory, markdown, nil)
		return
	}

	selected := make([]string, len(input.Categories))
	for i, n := range input.Categories {
		selected[i] = categories[n].Name
	}

	if stringInArray(category, selected) == -1 {
		input.Categories = append(input.Categories, index)
	} else {
		input.Categories = append(input.Categories[:index], input.Categories[index+1:]...)
	}

	reply := replyCategories
	for _, index := range input.Categories {
		reply += "\n" + names[index]
		buttons[index] = strAddMinus + names[index]
	}

	m.tgClient.SendMessage(msg.Chat.ID, reply,
		markdown, createAddKeyboards(buttons, reply != replyCategories))
}

func (m *manager) saveArticle(msg *tg.Message, input *db.StateInput) {
	if len(input.Categories) == 0 {
		m.tgClient.SendMessage(msg.Chat.ID, replyNoChosen, markdown, nil)
		return
	}

	err := m.dbClient.SaveArticle(input)
	if err != nil {
		m.SendError(msg.Chat.ID, err)
		return
	}

	input = nil
	m.tgClient.SendMessage(msg.Chat.ID, replyDone, markdown, defaultKeyboard)
}
