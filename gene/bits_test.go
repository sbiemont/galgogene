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

		Convey("new bits from bytes", func() {
			bits := NewBitsFromBytes([]byte{0x00, 0xFF, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80})
			So(bits, ShouldResemble, newBits([]uint8{
				0, 0, 0, 0, 0, 0, 0, 0,
				1, 1, 1, 1, 1, 1, 1, 1,
				0, 0, 0, 0, 0, 0, 0, 1,
				0, 0, 0, 0, 0, 0, 1, 0,
				0, 0, 0, 0, 0, 1, 0, 0,
				0, 0, 0, 0, 1, 0, 0, 0,
				0, 0, 0, 1, 0, 0, 0, 0,
				0, 0, 1, 0, 0, 0, 0, 0,
				0, 1, 0, 0, 0, 0, 0, 0,
				1, 0, 0, 0, 0, 0, 0, 0,
			}))
		})

		Convey("group", func() {
			Convey("when error nb bits", func() {
				bits := newBits([]uint8{0, 0, 0})
				by, err := bits.GroupBitsBy(2)
				So(err, ShouldBeError, "cannot group, total nb of bits (3) should be modulo 2")
				So(by, ShouldBeNil)
			})

			Convey("when error n>8", func() {
				by, err := Bits{}.GroupBitsBy(16)
				So(err, ShouldBeError, "cannot group, n > 8")
				So(by, ShouldBeNil)
			})

			Convey("when ok", func() {
				bits := newBits([]uint8{0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0})

				Convey("group by 2", func() {
					by, err := bits.GroupBitsBy(2)
					So(err, ShouldBeNil)
					So(by, ShouldResemble, []uint8{0b00, 0b01, 0b10, 0b11, 0b11, 0b10, 0b01, 0b00})
				})

				Convey("group by 4", func() {
					by, err := bits.GroupBitsBy(4)
					So(err, ShouldBeNil)
					So(by, ShouldResemble, []uint8{0b0001, 0b1011, 0b1110, 0b0100})
				})

				Convey("group by 8", func() {
					by, err := bits.GroupBitsBy(8)
					So(err, ShouldBeNil)
					So(by, ShouldResemble, []uint8{0b00011011, 0b11100100})
				})
			})
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
		})
	})
}
