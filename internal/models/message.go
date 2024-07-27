package models

import "time"

type Message struct {
	Uuid string
	Message string
	Count int8
	CreatedAt time.Time
}