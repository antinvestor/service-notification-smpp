package models

import (
	"github.com/pitabwire/frame"
)

// Template Table holds the test models
type Template struct {
	frame.BaseModel
	Name string `gorm:"type:varchar(255)"`
}
