package stu

import (
	"fmt"
	"strings"

	"github.com/researchlab/hybrid/examples/basic/lib/model"
	"github.com/researchlab/hybrid/orm"
)

// Service  student service struct
type Service struct {
	//DB orm.DBService `inject:"DB"`
	DB orm.DBService `inject:"DB"`
}

// SayHi sayhi for the given person by name
func (p *Service) SayHi(name string) (string, error) {
	one := &model.Stu{}
	if p.DB.GetDB().Where("name = ?", strings.TrimSpace(name)).Find(one).RecordNotFound() {
		return fmt.Sprintf("no one named %v found.", name), nil
	}
	res := fmt.Sprintf("Hi, Mr.%v.\n", one.Name)
	if one.Sex == "female" {
		res = fmt.Sprintf("Hi, Mrs.%v.\n", one.Name)
	}
	return res, nil
}
