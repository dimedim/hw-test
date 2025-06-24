package models

import "time"

type Notification struct {
	EventID  string
	Title    string
	StartsAt time.Time
	UserID   string
}
