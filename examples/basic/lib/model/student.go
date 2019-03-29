package model

import (
	"github.com/jinzhu/gorm"
	"github.com/researchlab/hybrid/orm"
)

// StuDesc student entity descriptor define
func StuDesc() *orm.ModelDescriptor {
	return &orm.ModelDescriptor{
		Type: &Stu{},
		New: func() interface{} {
			return &Stu{}
		},
		NewSlice: func() interface{} {
			return &[]Stu{}
		},
	}
}

// Stu student entity define
type Stu struct {
	gorm.Model
	Name    string
	Address string
	Sex     string
}
