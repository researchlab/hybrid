package orm_test

import (
	"fmt"
	"testing"
	"time"

	gorm "github.com/jinzhu/gorm"
	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/orm"
	"github.com/researchlab/hybrid/orm/dialects/mysql"
)

type A struct {
	gorm.Model
	Name string
	Bs   []*B
}

type B struct {
	gorm.Model
	AID  uint   `gorm:"index"`
	Name string `gorm:"unique"`
	Cs   []*C
}

type C struct {
	gorm.Model
	BID  uint `gorm:"index"`
	Name string
	Ds   []*D
}

type D struct {
	gorm.Model
	CID  uint `gorm:"index"`
	Name string
}

func aDesc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &A{},
		New: func() interface{} {
			return &A{}
		},
		NewSlice: func() interface{} {
			return &[]A{}
		},
	}
}

func bDesc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &B{},
		New: func() interface{} {
			return &B{}
		},
		NewSlice: func() interface{} {
			return &[]B{}
		},
	}
}

func cDesc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &C{},
		New: func() interface{} {
			return &C{}
		},
		NewSlice: func() interface{} {
			return &[]C{}
		},
	}
}

func dDesc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &D{},
		New: func() interface{} {
			return &D{}
		},
		NewSlice: func() interface{} {
			return &[]D{}
		},
	}
}

type Models struct {
	ModelRegistry orm.ModelRegistry `inject:"DB"`
}

func (p *Models) AfterNew() {
	p.ModelRegistry.Put("A", aDesc())
	p.ModelRegistry.Put("B", bDesc())
	p.ModelRegistry.Put("C", cDesc())
	p.ModelRegistry.Put("D", dDesc())
}

// Setup
type Setup struct {
	DB            orm.DBService     `inject:"DB"`
	ModelRegistry orm.ModelRegistry `inject:"DB"`
}

func (p *Setup) Init() error {

	p.DB.GetDB().LogMode(false)
	for v := range p.ModelRegistry.Models() {
		p.DB.GetDB().AutoMigrate(v.Type)
	}

	return nil
}

type inFunc func() error

type outFunc func() error

func TestCallback(t *testing.T) {
	container := brick.NewContainer()
	container.Add(&brick.JSONConfigService{}, "config", brick.FactoryFunc(func() interface{} {
		return brick.NewJSONConfigService("config/config.json")
	}))
	container.Add(&Models{}, "Models", nil)
	container.Add(&mysql.MySQLService{}, "DB", nil)
	container.Add(&Setup{}, "Setup", nil)
	container.Build()
	defer container.Dispose()

	dbSvc, ok := container.GetByName("DB").(*mysql.MySQLService)
	if !ok {
		t.Fatal("Cloud not found DB")
	}

	testCases := []struct {
		name string
		in   inFunc
		out  outFunc
	}{
		{
			name: "delete multi cascade",
			in: func() error {
				a := &A{Name: "a", Bs: []*B{}}
				b1 := &B{Name: "b1", Cs: []*C{}}
				b2 := &B{Name: "b2", Cs: []*C{}}
				c1 := &C{Name: "c1"}
				c2 := &C{Name: "c2"}
				d1 := &D{Name: "d1"}
				d2 := &D{Name: "d2"}
				a.Bs = append(a.Bs, b1)
				a.Bs = append(a.Bs, b2)
				b1.Cs = append(b1.Cs, c1)
				b2.Cs = append(b2.Cs, c2)
				c1.Ds = append(c1.Ds, d1)
				c2.Ds = append(c2.Ds, d2)
				if err := dbSvc.GetDB().Create(a).Error; err != nil {
					return err
				}

				if err := dbSvc.GetDB().Unscoped().Delete(a).Error; err != nil {
					return err
				}
				return nil
			},
			out: func() error { return nil },
		},
		{
			name: "update delete multi cascade",
			in: func() error {
				a := &A{Name: "a", Bs: []*B{}}
				b1 := &B{Name: "bu1", Cs: []*C{}}
				b2 := &B{Name: "bu2", Cs: []*C{}}
				c1 := &C{Name: "cu1"}
				c2 := &C{Name: "cu2"}
				d1 := &D{Name: "du1"}
				d2 := &D{Name: "du2"}
				a.Bs = append(a.Bs, b1)
				a.Bs = append(a.Bs, b2)
				b1.Cs = append(b1.Cs, c1)
				b2.Cs = append(b2.Cs, c2)
				c1.Ds = append(c1.Ds, d1)
				c2.Ds = append(c2.Ds, d2)
				if err := dbSvc.GetDB().Create(a).Error; err != nil {
					return err
				}

				t.Logf("%v", a.ID)
				ra := &A{}
				if dbSvc.GetDB().Preload("Bs").Preload("Bs.Cs").First(ra, a.ID).RecordNotFound() {
					return fmt.Errorf("could not found data %d", a.ID)
				}
				t.Logf("%+v", ra)

				ra.Bs = ra.Bs[1:]
				var empty time.Time
				ra.Bs[0].CreatedAt = empty
				ra.Bs = append(ra.Bs, &B{Name: "bu3", Cs: []*C{}})
				fmt.Println("update")
				if err := dbSvc.GetDB().Save(ra).Error; err != nil {
					return err
				}

				if dbSvc.GetDB().Preload("Bs").Preload("Bs.Cs").First(ra, a.ID).RecordNotFound() {
					return fmt.Errorf("could not found data %d", a.ID)
				}
				t.Logf("%+v", ra)
				return nil
			},
			out: func() error { return nil },
		},
	}

	for _, testCase := range testCases {
		if err := testCase.in(); err != nil {
			t.Errorf("case %s execution fail: %s", testCase.name, err.Error())
		}
		if err := testCase.out(); err != nil {
			t.Errorf("case %s valid fail: %s", testCase.name, err.Error())
		}
	}
}
