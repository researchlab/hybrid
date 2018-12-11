package gorm

import "time"

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      gorm.Model
//    }
type Model struct {
	ID        uint   `gorm:"primary_key"`
	CreatedBy string `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedBy string `gorm:"size:50"`
	UpdatedAt time.Time
	DeletedBy string     `gorm:"size:50"`
	DeletedAt *time.Time `sql:"index"`
}
