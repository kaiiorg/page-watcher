package models

import (
	"time"
)

type Page struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string
	Text string
	Diff string
}
