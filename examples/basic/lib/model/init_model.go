package model

import "github.com/researchlab/hybrid/orm"

// Models db models
type Models struct {
	ModelRegistry orm.ModelRegistry `inject:"DB"`
}

// AfterNew registry db model
func (p *Models) AfterNew() {
	p.ModelRegistry.Put("Stu", StuDesc())
}
