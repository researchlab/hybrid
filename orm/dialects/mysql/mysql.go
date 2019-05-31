package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // import mysql driven
	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/orm"
)

// Service mysql service provider
type Service struct {
	orm.ModelRegistryImpl
	db     *gorm.DB
	Config brick.Config `inject:"config"`
}

// Init mysql service init
func (p *Service) Init() error {
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

	//timeout, err := time.ParseDuration(p.Config.GetMapString("db", "connMaxLifetime", "2h"))
	//if err != nil {
	//	timeout = 2 * time.Hour
	//}
	//db.DB().SetConnMaxLifetime(3 * time.Second)

	//fix the bug: [mysql] 2016/10/11 09:17:16 packets.go:33: unexpected EOF
	db.DB().SetMaxIdleConns(0)

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

// Dispose close mysql conn
func (p *Service) Dispose() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// GetDB return db conn
func (p *Service) GetDB() *gorm.DB {
	return p.db
}
