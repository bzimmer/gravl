package activity

import "time"

// A Handle provides a minimal set of metadata about an entity
type Handle struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Source string    `json:"source"`
	Date   time.Time `json:"date"`
}

// A Named entity supports returning a Handle for an entity
type Named interface {
	Handle() *Handle
}
