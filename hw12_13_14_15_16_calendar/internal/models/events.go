package models

import "time"

type Event struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description,omitempty" db:"description"`

	StartsAt     time.Time     `json:"starts_at" db:"starts_at"`
	EndsAt       time.Time     `json:"ends_at,omitzero" db:"ends_at"`
	NotifyOffset time.Duration `json:"notify_offset,omitempty" db:"notify_offset"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitzero" db:"updated_at"`
}
