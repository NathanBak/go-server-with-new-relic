package widget

import (
	"time"

	"github.com/google/uuid"
)

// A Widget is an example struct
type Widget struct {
	Name  string    `json:"name"`
	ID    uuid.UUID `json:"id"`
	Color Color     `json:"color"`
	Time  time.Time `json:"lastUpdated"`
}

// New creates and returns a new Widget with a unique ID.
func New(name string, color Color) Widget {
	return Widget{name, uuid.New(), color, time.Now()}
}
