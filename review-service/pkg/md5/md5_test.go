package md5

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSum(t *testing.T) {
	convey.Convey("test1", t, func() {
		var (
			initialstring = "mysteriousx"
			md5string     = "23b63e56b868e7d6f85285f76a6d6ec6"
		)
		md5res := Sum([]byte(initialstring))
		convey.So(md5res, convey.ShouldEqual, md5string)
	})

	convey.Convey("test1", t, func() {
		var (
			initialstring = "mysteriousx"
			md5string     = "b868e7d6f85285f7"
		)
		md5res := Sum([]byte(initialstring))
		convey.So(md5res, convey.ShouldNotEqual, md5string)
	})
}
