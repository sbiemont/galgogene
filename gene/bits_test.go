package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func newBits(bits []uint8) Bits {
	return Bits{
		Raw:      bits,
		MaxValue: DefaultMaxValue,
	}
}

func TestBits(t *testing.T) {
	Convey("bits", t, func() {
		Convey("new bits", func() {
			bits := NewBits(8, DefaultMaxValue)
			So(bits, ShouldResemble, newBits([]uint8{0, 0, 0, 0, 0, 0, 0, 0}))
		})

		Convey("new bits random", func() {
			n := 1000
			bits := NewBitsRandom(n, DefaultMaxValue)
			for i := 0; i < n; i++ {
				So(bits.Raw[i], ShouldBeBetweenOrEqual, 0, 1)
			}
		})

		Convey("transform", func() {
			Convey("when max value: 1", func() {
				bits := Bits{MaxValue: 1}
				So(bits.modulo(0), ShouldEqual, 0)
				So(bits.modulo(1), ShouldEqual, 1)
				So(bits.modulo(2), ShouldEqual, 0)
				So(bits.modulo(3), ShouldEqual, 1)
			})

			Convey("when max value: 255", func() {
				bits := Bits{MaxValue: 255}
				So(bits.modulo(0), ShouldEqual, 0)
				So(bits.modulo(1), ShouldEqual, 1)
				So(bits.modulo(254), ShouldEqual, 254)
				So(bits.modulo(255), ShouldEqual, 255)
			})

			Convey("when invert: 1", func() {
				bits := Bits{
					MaxValue: 1,
					Raw:      []uint8{0, 1},
				}
				So(bits.Invert(0), ShouldEqual, 1)
				So(bits.Invert(1), ShouldEqual, 0)
			})

			Convey("when invert: 255", func() {
				bits := Bits{
					MaxValue: 255,
					Raw:      []uint8{0, 1, 254, 255},
				}
				So(bits.Invert(0), ShouldEqual, 255)
				So(bits.Invert(1), ShouldEqual, 254)
				So(bits.Invert(2), ShouldEqual, 1)
				So(bits.Invert(3), ShouldEqual, 0)
			})

			Convey("when bytes", func() {
				bits := Bits{
					Raw: []uint8{0, 1, 254, 255},
				}
				So(bits.Bytes(), ShouldResemble, []byte{0, 1, 254, 255})
			})
		})
	})
}
