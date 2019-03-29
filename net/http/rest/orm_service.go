package rest

import (
	"fmt"
	"strings"

	"github.com/researchlab/hybrid/orm"
)

type ormService struct {
	db            orm.DBService
	modelRegistry orm.ModelRegistry
	watcher       *watcher
}

func newService(db orm.DBService, modelRegistry orm.ModelRegistry) *ormService {
	return &ormService{db: db, modelRegistry: modelRegistry, watcher: newWatcher()}
}

func (p *ormService) Get(class string, id interface{}, ass string) (interface{}, error) {
	md := p.modelRegistry.Get(class)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", class)
	}

	data := md.New()
	d := p.db.GetDB()
	for _, as := range strings.Split(ass, ",") {
		d = d.Preload(as)
	}

	if d.First(data, id).RecordNotFound() {
		return nil, nil
	}

	return data, nil
}

func (p *ormService) List(class string, selectFields []string, where string, whereValues []interface{}, order string, page int, pageSize int) (map[string]interface{}, error) {
	md := p.modelRegistry.Get(class)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", class)

	}

	//count
	d := p.db.GetDB().Model(md.NewSlice())
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
		d = p.db.GetDB().Where(where, whereValues...)
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

func (p *ormService) Create(class string, data interface{}) error {
	if err := p.db.GetDB().Create(data).Error; err != nil {
		return err
	}

	//p.watcher.NotifyCreate(class(class), data)
	return nil
}

type classType string

func (p *ormService) Remove(className string, id interface{}) (interface{}, error) {
	md := p.modelRegistry.Get(className)
	if md == nil {
		return nil, fmt.Errorf("class %s isn't exists", className)
	}
	data := md.New()
	if err := p.db.GetDB().First(data, id).Error; err != nil {
		return nil, err
	}

	if err := p.db.GetDB().Unscoped().Delete(data).Error; err != nil {
		return nil, err
	}

	p.watcher.NotifyDelete(class(className), id)
	return data, nil
}

func (p *ormService) Update(className string, data interface{}) error {
	if err := p.db.GetDB().Save(data).Error; err != nil {
		return err
	}
	p.watcher.NotifyUpdate(class(className), data)
	return nil
}

func (p *ormService) WatchObject(wk int64, timeout int, className string, id interface{}, associations string) (int64, chan *WatchEvent) {
	return p.watcher.Wait(wk, timeout, class(className), id)
}
