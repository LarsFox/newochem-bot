package db

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type client struct {
	db *gorm.DB
}

// Client is a psql gorm connecting client
type Client interface {
	CloseDB()
	GetCategories() ([]*Category, error)
	GetWorkers() ([]*Worker, error)
	GetState(user string) (*State, *StateInput, error)
	CreateState(user string) (*State, error)
	SetState(s *State, input *StateInput) error
}

// Closes DB from main
func (c *client) CloseDB() {
	c.db.Close()
}

// NewClient opens a new mysql client and returns it
func NewClient(dbString string) (Client, error) {
	db, err := gorm.Open("mysql", dbString)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &client{db: db}, nil
}

func (c *client) GetCategories() ([]*Category, error) {
	categories := []*Category{}
	c.db.Where("active = 1").Find(&categories)
	return categories, nil
}

func (c *client) GetWorkers() ([]*Worker, error) {
	workers := []*Worker{}
	c.db.Where("active = 1").Find(&workers)
	return workers, nil
}

func (c *client) GetState(user string) (*State, *StateInput, error) {
	s := &State{}
	err := c.db.Where("user = ?", user).First(s).Error
	if err != nil {
		return nil, nil, err
	}

	input := &StateInput{}
	err = json.Unmarshal([]byte(s.Input), input)
	if err != nil {
		return nil, nil, err
	}
	return s, input, nil
}

// TODO: error
func (c *client) CreateState(user string) (*State, error) {
	s := &State{User: user, Input: "{}"}
	c.db.Save(s)
	return s, nil
}

func (c *client) SetState(s *State, input *StateInput) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}
	s.Input = string(b)
	c.db.Save(s)
	return nil
}
