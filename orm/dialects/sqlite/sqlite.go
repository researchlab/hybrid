package sqlite

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/researchlab/hybrid/brick"
	"github.com/researchlab/hybrid/orm"
)

// SQLite3Service
type SQLite3Service struct {
	orm.ModelRegistryImpl
	orm.DeleteCallback
	orm.UpdateCallback
	db     *gorm.DB
	Config brick.Config `inject:"config"`
}

func (p *SQLite3Service) Init() error {
	p.DeleteCallback.Register(p)
	p.UpdateCallback.Register(p)

	filename := p.Config.GetMapString("db", "url", "../data/nerv.db")
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			return fmt.Errorf("create dir %s failed. %s", dir, err.Error())
		}
	}

	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		return err
	}

	db.LogMode(false)
	for v := range p.Models() {
		db.AutoMigrate(v.Type)
	}
	if db.Error != nil {
		return db.Error
	}
	p.db = db
	return nil
}

func (p *SQLite3Service) Dispose() error {
	return p.db.Close()
}

func (p *SQLite3Service) GetDB() *gorm.DB {
	return p.db
}
