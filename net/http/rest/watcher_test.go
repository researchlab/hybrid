package rest

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWait(t *testing.T) {
	watcher := newWatcher()
	_, ch := watcher.Wait(0, 10, "Alert", 1)

	//	e.OrmEvent
	Convey("wait object", t, func() {

		Convey("create object", func() {
			o := &struct{}{}
			go func() {
				watcher.NotifyCreate("Alert", o)
			}()

			e := <-ch
			So(e.OrmEvent.EventType, ShouldEqual, orm_create)
			So(e.OrmEvent.Data, ShouldEqual, o)
		})

		Convey("update object", func() {
			o := &struct{}{}
			go func() {
				watcher.NotifyUpdate("Alert", o)
			}()
			e := <-ch
			So(e.OrmEvent.EventType, ShouldEqual, orm_update)
			So(e.OrmEvent.Data, ShouldEqual, o)
		})

		Convey("delete object", func() {
			id := 1
			go func() {
				watcher.NotifyDelete("Alert", id)
			}()
			e := <-ch
			So(e.OrmEvent.EventType, ShouldEqual, orm_delete)
			So(e.OrmEvent.ID, ShouldEqual, id)
		})
	})

}

func TestNotifyUpdate(t *testing.T) {

}

func TestNotifyDelete(t *testing.T) {

}

func TestNotifyWatch(t *testing.T) {

}
