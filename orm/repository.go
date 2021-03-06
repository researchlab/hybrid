package orm

import (
	"github.com/jinzhu/gorm"
)

// Repository manage all objects in db
type Repository interface {
	Get(class string, id interface{}, ass string) (interface{}, error)
	List(class string, selectFields []string, where string, whereValues []interface{}, order string, page int, pageSize int) (map[string]interface{}, error)
	Create(class string, data interface{}) error
	Remove(className string, id interface{}) (interface{}, error)
	Update(className string, data interface{}) error
	GetDB() *gorm.DB
}
