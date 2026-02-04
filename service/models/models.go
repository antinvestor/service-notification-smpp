package models

import (
	"github.com/pitabwire/frame/data"
)

// Template Table holds the test models
type Template struct {
	data.BaseModel
	Name string `gorm:"type:varchar(255)"`
}
