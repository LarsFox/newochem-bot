package tg

import "encoding/json"

// Message is a Telegram message
type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	Chat      *Chat  `json:"chat"`
}

// Chat is a Telegram chat
type Chat struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Update https://core.telegram.org/bots/api#update
type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

// Response is a Telegram response
type Response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description json.RawMessage `json:"description"` // only if not OK
	ErrorCode   int             `json:"error_code"`
}

// ReplyKeyboardMarkup is a bot Keyboard
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`   // optional
	OneTimeKeyboard bool               `json:"one_time_keyboard"` // optional
	Selective       bool               `json:"selective"`         // optional
}

// KeyboardButton one key in ReplyKeyboardMarkup.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

/*
// ReplyKeyboardHide allows the Bot to hide a custom keyboard.
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"` // optional
}

// ReplyKeyboardRemove allows the Bot to hide a custom keyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}
*/
