package model

import uuid "github.com/satori/go.uuid"

type Task struct {
	ID            uuid.UUID         `json:"id"`
	Status        int               `json:"status"`
	ContentLength int               `json:"content_length"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"-"`
}
