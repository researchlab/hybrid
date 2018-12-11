package orm

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
)

// DeleteCallback
type DeleteCallback struct {
	models ModelRegistry
}

func (p *DeleteCallback) Register(models ModelRegistry) {
	p.models = models
	gorm.DefaultCallback.Delete().Replace("gorm:begin_transaction", p.beginTransactionCallback)
	gorm.DefaultCallback.Delete().Replace("gorm:before_delete", p.beforeDeleteCallback)
	gorm.DefaultCallback.Delete().After("gorm:before_delete").Replace("hybrid:before_delete_associations", p.beforeDeleteAssociationsCallback)
	gorm.DefaultCallback.Delete().Replace("gorm:delete", p.deleteCallback)
	gorm.DefaultCallback.Delete().After("gorm:delete").Replace("hybrid:after_delete_associations", p.afterDeleteAssociationsCallback)
	gorm.DefaultCallback.Delete().Replace("gorm:after_delete", p.afterDeleteCallback)
	gorm.DefaultCallback.Delete().Replace("gorm:commit_or_rollback_transaction", p.commitOrRollbackTransactionCallback)
}

func (p *DeleteCallback) beginTransactionCallback(scope *gorm.Scope) {
	scope.Begin()
}

// beforeDeleteCallback will invoke `BeforeDelete` method before deleting
func (p *DeleteCallback) beforeDeleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		scope.CallMethod("BeforeDelete")
	}
}

// DeleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func (p *DeleteCallback) deleteCallback(scope *gorm.Scope) {
	//log.LogCodeLine()
	//fmt.Printf("deleteCallback: %s\n", scope.TableName())
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		if !scope.Search.Unscoped && scope.HasColumn("DeletedAt") {

			ctx := scope.GetContext()
			v := ctx.Value(gorm.ContextCurrentUser())
			var user string
			if v != nil {
				user = v.(string)
			}
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET deleted_at=%v,deleted_by=%v %v%v",
				scope.QuotedTableName(),
				scope.AddToVars(gorm.NowFunc()),
				scope.AddToVars(user),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			sql := fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption))
			////fmt.Println(sql)
			scope.Raw(sql).Exec()
		}
	}
}

// afterDeleteCallback will invoke `AfterDelete` method after deleting
func (p *DeleteCallback) afterDeleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterDelete")
	}
}

func (p *DeleteCallback) commitOrRollbackTransactionCallback(scope *gorm.Scope) {
	scope.CommitOrRollback()
}

//cascade deleting
func (p *DeleteCallback) beforeDeleteAssociationsCallback(scope *gorm.Scope) {
	//log.LogCodeLine()
	//TBD config  gorm:delete_associations
	//if !scope.shouldDeleteAssociations() {
	//	return
	//}
	//fmt.Printf("beforeDeleteAssociationsCallback: %s\n", scope.TableName())
	for _, field := range scope.Fields() {
		//fmt.Printf("beforeDeleteAssociationsCallback: %s.%s\n", scope.TableName(), field.Name)
		if relationship := field.Relationship; relationship != nil && relationship.Kind == "has_many" {
			//TBD:Now only support one foreign field and unit type
			//foreignValue := scope.IndirectValue().FieldByName("ID").Uint()
			//fmt.Printf("%s %+v\n", relationship.AssociationForeignFieldNames[0],scope.IndirectValue().Interface())
			foreignValue := scope.IndirectValue().FieldByName(relationship.AssociationForeignFieldNames[0]).Uint()
			sql := fmt.Sprintf("%s = ?", relationship.ForeignDBNames[0])
			fieldType := field.Field.Type()
			elem := fieldType.Elem()
			class := ""
			if elem.Kind() == reflect.Ptr {
				class = elem.Elem().Name()
			} else {
				class = elem.Name()
			}
			//fmt.Printf("beforeDeleteAssociationsCallback: %s %s %s %d\n", scope.TableName(), class, sql, foreignValue)
			children := p.models.Get(class).NewSlice()

			if !scope.DB().Where(sql, foreignValue).Find(children).RecordNotFound() {
				cv := reflect.ValueOf(children).Elem()
				len := reflect.ValueOf(cv.Interface()).Len()

				for i := 0; i < len; i++ {
					child := cv.Index(i).Interface()
					if err := scope.DB().Unscoped().Delete(child).Error; err != nil {
						scope.Err(err)
					}
				}
			}
			// if err := scope.DB().Unscoped().Delete(p.models.Get(class).Type, sql, foreignValue).Error; err != nil {
			// 	//log.LogCodeLine()
			// 	scope.Err(err)
			// }
		}
	}
}

func (p *DeleteCallback) afterDeleteAssociationsCallback(scope *gorm.Scope) {
	//log.LogCodeLine()
}
