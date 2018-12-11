package model

import "github.com/researchlab/hybrid/orm"

type Models struct {
	ModelRegistry orm.ModelRegistry `inject:"DB"`
}

func (p *Models) AfterNew() {
	p.ModelRegistry.Put("Stu", StuDesc())
}
