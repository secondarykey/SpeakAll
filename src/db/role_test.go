package db

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRole(t *testing.T) {

	var err error
	// Only pass t into top-level Convey calls
	Convey("database listen", t, func() {

		err = Listen("../../data/db/test.db")
		So(err, ShouldBeNil)

		Convey("create role table", func() {
		})

		Convey("insert role table", func() {
		})
	})

}