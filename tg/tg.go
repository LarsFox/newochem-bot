package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

// Client works with Telegram Bot API.
type Client interface {
	GetMessagesChan() (<-chan *Message, error)
	SendAudio(chatID int, audio []byte) (*Response, error)
	SendMessage(chatID int, text, parse string, keyboard *ReplyKeyboardMarkup) (*Response, error)
	SendPhoto(chatID int, photo []byte, caption string) (*Response, error)
	SendDocument(chatID int, file []byte, fileName, caption string) (*Response, error)
}

type client struct {
	token     string
	lastUpdID int
}

// NewClient returns a new client to work with Telegram Bot API.
func NewClient(token string) Client {
	return &client{token: token}
}

// GetMessagesChan returns a message channel
// TODO: errors
func (c *client) GetMessagesChan() (<-chan *Message, error) {
	msgChan := make(chan *Message)

	go func() {
		for {
			uri := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=100&offset=%d",
				c.token, c.lastUpdID+1)
			resp, err := c.makeGETRequest(uri)
			if err != nil {
				continue
			}

			var updates []*Update
			err = json.Unmarshal(resp.Result, &updates)
			if err != nil {
				continue
			}

			for _, upd := range updates {
				c.lastUpdID = upd.UpdateID
				msgChan <- upd.Message
			}
		}
	}()

	return msgChan, nil
}

func (c *client) makeGETRequest(uri string) (*Response, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) makeTelegramRequest(methodName, fileKey, fileName string, file []byte,
	fields map[string]string) (*Response, error) {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/"+methodName, c.token)
	resp, err := c.makePOSTRequest(uri, fileKey, fileName, file, fields)

	if !resp.Ok {
		err = fmt.Errorf("Tg error %d", resp.ErrorCode)
		return nil, err
	}

	return resp, nil
}

func (c *client) makePOSTRequest(uri, fileKey, fileName string, file []byte,
	fields map[string]string) (*Response, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	// Do nothing with empty file
	if len(file) > 0 {
		// TODO: Buffer -> []byte -> NewBuffer([]byte)?
		buf := bytes.NewBuffer(file)

		fw, err := w.CreateFormFile(fileKey, fileName)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fw, buf); err != nil {
			return nil, err
		}
	}

	for k, v := range fields {
		w.WriteField(k, v)
	}

	w.Close()

	req, err := http.NewRequest("POST", uri, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &Response{}
	err = json.Unmarshal(body, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Sends Audio
func (c *client) SendAudio(chatID int, audio []byte) (*Response, error) {
	return c.makeTelegramRequest("sendAudio", "audio", "audio.mp3", audio, map[string]string{
		"chat_id": strconv.Itoa(chatID),
	})
}

// Sends a message
func (c *client) SendMessage(chatID int, text, parse string,
	keyboard *ReplyKeyboardMarkup) (*Response, error) {
	var key string
	if keyboard != nil {
		data, err := json.Marshal(keyboard)
		if err != nil {
			return nil, err
		}

		key = string(data)
	}

	return c.makeTelegramRequest("sendMessage", "", "", []byte{}, map[string]string{
		"chat_id":      strconv.Itoa(chatID),
		"text":         text,
		"parse_mode":   parse,
		"reply_markup": key,
		// disable_web_page_preview	Boolean
		// disable_notification	Boolean
	})
}

// Sends a photo (image)
func (c *client) SendPhoto(chatID int, photo []byte, caption string) (*Response, error) {
	return c.makeTelegramRequest("sendPhoto", "photo", "photo.png", photo, map[string]string{
		"chat_id": strconv.Itoa(chatID),
		"caption": caption,
		// disable_notification	Boolean
	})
}

// Sends a document
func (c *client) SendDocument(chatID int, file []byte, fileName, caption string) (*Response, error) {
	return c.makeTelegramRequest("sendDocument", "document", fileName, file, map[string]string{
		"chat_id": strconv.Itoa(chatID),
		"caption": caption,
		// disable_notification	Boolean
	})
}
