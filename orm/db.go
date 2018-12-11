package orm

import (
	"github.com/jinzhu/gorm"
)

// DBService
type DBService interface {
	GetDB() *gorm.DB
}
