package mysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/orm"
)

// MySQLService
type MySQLService struct {
	orm.ModelRegistryImpl
	db     *gorm.DB
	Config brick.Config `inject:"config"`
}

func (p *MySQLService) Init() error {
	passwd := p.Config.GetMapString("db", "password", "root")
	url := fmt.Sprintf(
		"%s:%s@%s",
		p.Config.GetMapString("db", "user", "root"),
		passwd,
		p.Config.GetMapString("db", "url"),
	)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return err
	}

	timeout, err := time.ParseDuration(p.Config.GetMapString("db", "connMaxLifetime", "2h"))
	if err != nil {
		timeout = 2 * time.Hour
	}
	db.DB().SetConnMaxLifetime(timeout)

	log := p.Config.GetMapBool("db", "log", false)
	db.LogMode(log)
	for v := range p.Models() {
		if err := db.AutoMigrate(v.Type).Error; err != nil {
			return err
		}
	}

	p.db = db
	return nil
}

func (p *MySQLService) Dispose() error {
	if p.db != nil {
		return p.db.Close()
	} else {
		return nil
	}
}

func (p *MySQLService) GetDB() *gorm.DB {
	return p.db
}
