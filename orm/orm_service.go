package orm

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// InstanceNotFoundError xxx
type InstanceNotFoundError struct {
	msg string
}

func newInstanceNotFoundError(msg string) *InstanceNotFoundError {
	return &InstanceNotFoundError{msg: msg}
}

func (p *InstanceNotFoundError) Error() string {
	return p.msg
}

//Service orm service struct
type Service struct {
	DB            DBService     `inject:"DB"`
	ModelRegistry ModelRegistry `inject:"DB"`
}

//GetDB return db conn
func (p *Service) GetDB() *gorm.DB {
	return p.DB.GetDB()
}

//Get query db model with where conditions
func (p *Service) Get(class string, id interface{}, ass string) (interface{}, error) {
	md := p.ModelRegistry.Get(class)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", class)
	}

	data := md.New()
	d := p.DB.GetDB()
	for _, as := range strings.Split(ass, ",") {
		d = d.Preload(as)
	}

	if d.First(data, id).RecordNotFound() {
		return nil, newInstanceNotFoundError(fmt.Sprintf("could not found instance %v of class %s", id, class))
	}

	return data, nil
}

//List  query db model list with where conditions
func (p *Service) List(class string, selectFields []string, where string, whereValues []interface{}, order string, page int, pageSize int) (map[string]interface{}, error) {
	md := p.ModelRegistry.Get(class)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", class)

	}

	//count
	d := p.DB.GetDB().Model(md.NewSlice())
	var count int64
	if where != "" {
		d = d.Where(where, whereValues...)
	}
	if len(selectFields) > 0 {
		d = d.Select(selectFields)
	}
	if err := d.Count(&count).Error; err != nil {
		return nil, err
	}
	//order page
	var pageCount, limit int
	if pageSize > 0 {
		limit = pageSize
	} else {
		limit = 10
	}
	pageCount = (int(count) + limit - 1) / limit

	if page < 0 {
		page = 0
	}

	if page >= pageCount {
		page = pageCount - 1
	}

	if where != "" {
		d = p.DB.GetDB().Where(where, whereValues...)
	}
	if len(selectFields) > 0 {
		d = d.Select(selectFields)
	}
	if order != "" {
		d = d.Order(order)
	}
	data := md.NewSlice()
	if d.Offset(page * limit).Limit(limit).Find(data).RecordNotFound() {
		return nil, nil
	}

	return map[string]interface{}{"data": data, "page": page, "pageSize": limit, "pageCount": pageCount}, nil
}

//Create create data
func (p *Service) Create(class string, data interface{}) error {
	if err := p.DB.GetDB().Create(data).Error; err != nil {
		return err
	}

	//p.watcher.NotifyCreate(class(class), data)
	return nil
}

type classType string

//Remove remove db data
func (p *Service) Remove(className string, id interface{}) (interface{}, error) {
	md := p.ModelRegistry.Get(className)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", className)
	}
	data := md.New()
	if err := p.DB.GetDB().First(data, id).Error; err != nil {
		return nil, err
	}

	if err := p.DB.GetDB().Unscoped().Delete(data).Error; err != nil {
		return nil, err
	}

	//p.watcher.NotifyDelete(class(className), id)
	return data, nil
}

//Update update db data
func (p *Service) Update(className string, data interface{}) error {
	if err := p.DB.GetDB().Save(data).Error; err != nil {
		return err
	}
	//p.watcher.NotifyUpdate(class(className), data)
	return nil
}

// func (p *Service) WatchObject(wk int64, timeout int, className string, id interface{}, associations string) (int64, chan *WatchEvent) {
// 	return p.watcher.Wait(wk, timeout, class(className), id)
// }
