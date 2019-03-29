package orm

import (
	"github.com/jinzhu/gorm"
)

// DBService DB Service interface
type DBService interface {
	GetDB() *gorm.DB
}
