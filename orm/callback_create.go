package orm

import (
	"github.com/jinzhu/gorm"
)

// CreateCallback
type CreateCallback struct {
	models ModelRegistry
}

func (p *CreateCallback) Register(models ModelRegistry) {
	p.models = models
	gorm.DefaultCallback.Create().Replace("gorm:update_time_stamp", p.updateTimeStampForCreateCallback)
}

func (p *CreateCallback) updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		//fmt.Printf("updateTimeStampForCreateCallback: %s\n", scope.GetModelStruct().ModelType.Name())
		now := NowFunc()
		scope.SetColumn("CreatedAt", now)
		scope.SetColumn("UpdatedAt", now)

		ctx := scope.GetContext()
		v := ctx.Value(gorm.ContextCurrentUser())
		switch user := v.(type) {
		case string:
			if field, ok := scope.FieldByName("CreatedBy"); ok {
				field.Set(user)
			}
			if field, ok := scope.FieldByName("UpdatedBy"); ok {
				field.Set(user)
			}
		}
	}
}
