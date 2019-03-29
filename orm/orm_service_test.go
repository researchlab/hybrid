package orm_test

import (
	"github.com/jinzhu/gorm"

	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/orm"
	"github.com/researchlab/hybrid/orm/dialects/mysql"
)

var ormSvc *orm.Service

func init() {
	c := brick.NewContainer()
	c.Add(&brick.JSONConfigService{}, "config", brick.FactoryFunc(func() interface{} {
		return brick.NewJSONConfigService("config/config.json")
	}))

	c.Add(&orm.Service{}, "Service", nil)
	c.Add(&mysql.MySQLService{}, "DB", nil)
	c.Add(&TestModel{}, "TestEntity", nil)
	c.Build()
	ormSvc = c.GetByName("Service").(*orm.Service)
}

type TestEntity struct {
	gorm.Model
	Name string
}

type TestModel struct {
	ModelRegistry orm.ModelRegistry `inject:"DB"`
}

func (p *TestModel) AfterNew() {
	p.ModelRegistry.Put("TestEntity", p.desc())
}

func (p *TestModel) desc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &TestEntity{},
		New: func() interface{} {
			return &TestEntity{}
		},
		NewSlice: func() interface{} {
			return &[]TestEntity{}
		},
	}
}
