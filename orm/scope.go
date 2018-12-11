package orm

import "github.com/jinzhu/gorm"

func shouldSaveAssociations(scope *gorm.Scope) bool {
	if saveAssociations, ok := scope.Get("gorm:save_associations"); ok && !saveAssociations.(bool) {
		return false
	}
	return true && !scope.HasError()
}

func changeableField(scope *gorm.Scope, field *gorm.Field) bool {
	if selectAttrs := scope.SelectAttrs(); len(selectAttrs) > 0 {
		for _, attr := range selectAttrs {
			if field.Name == attr || field.DBName == attr {
				return true
			}
		}
		return false
	}

	for _, attr := range scope.OmitAttrs() {
		if field.Name == attr || field.DBName == attr {
			return false
		}
	}

	return true
}
