package core

import "github.com/larsfox/newochem-bot/tg"

const keysRow = 3

var (
	defaultKeyboard = &tg.ReplyKeyboardMarkup{
		// OneTimeKeyboard: true,
		ResizeKeyboard: true,
		Keyboard: [][]tg.KeyboardButton{
			{{Text: strAddArticle}},
		},
	}

	lastRow = []tg.KeyboardButton{tg.KeyboardButton{Text: strCancel}}
)

func createAddKeyboards(stringsArray []string, ready bool) *tg.ReplyKeyboardMarkup {
	keyboard := make([][]tg.KeyboardButton, 0, len(stringsArray))
	for _, chunk := range chunksString(stringsArray, keysRow) {
		row := []tg.KeyboardButton{}
		for _, name := range chunk {
			row = append(row, tg.KeyboardButton{Text: name})
		}
		keyboard = append(keyboard, row)
	}

	last := lastRow
	if ready {
		last = append(lastRow, tg.KeyboardButton{Text: strDone})
	}

	return &tg.ReplyKeyboardMarkup{
		// OneTimeKeyboard: true,
		ResizeKeyboard: true,
		Keyboard:       append(keyboard, last),
	}
}
