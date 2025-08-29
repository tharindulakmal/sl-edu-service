package models

type Grade struct {
	ID    int    `json:"id" db:"id"`
	Grade string `json:"grade" db:"grade"`
}
