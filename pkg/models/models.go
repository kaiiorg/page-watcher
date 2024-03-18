package models

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Name string
	Text string
	Diff string
}
