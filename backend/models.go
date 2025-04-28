package main

import (
	"time"

	"github.com/google/uuid"
)

type Reading struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Cards     []string  `json:"cards"`
	Timestamp time.Time `json:"timestamp"`
}

type TarotCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Keywords    string `json:"keywords"`
}

var (
	readingsDB = make(map[string][]Reading)
	cardsDB    = generateTarotCards()
)

func generateTarotCards() []TarotCard {
	return []TarotCard{
		{
			ID:          uuid.New().String(),
			Name:        "Шут",
			Description: "Символ новых начинаний...",
			ImageURL:    "fool.jpg",
			Keywords:    "начало, невинность",
		},
		// othershit
	}
}
