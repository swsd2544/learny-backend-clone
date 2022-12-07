package entity

import "time"

const (
	COMMON    = "common"
	RARE      = "rare"
	LEGENDARY = "legendary"
	MYSTIC    = "mystic"
)

type Character struct {
	ID        int64     `json:"id"`
	ImageURL  string    `json:"image_url"`
	Rarity    string    `json:"string"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"-"`
}
