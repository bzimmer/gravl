package activity

import "time"

// A Named provides a minimal set of metadata about an entity
type Named struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Source string    `json:"source"`
	Date   time.Time `json:"date"`
}

// A Namer returns a Named for an entity
type Namer interface {
	Named() *Named
}
