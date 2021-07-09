package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBits(t *testing.T) {
	Convey("bits", t, func() {
		Convey("new bits", func() {
			bits := NewBits(8)
			So(bits, ShouldResemble, Bits{0, 0, 0, 0, 0, 0, 0, 0})
		})

		Convey("new bits random", func() {
			n := 1000
			bits := NewBitsRandom(n)
			for i := 0; i < n; i++ {
				So(bits[i], ShouldBeBetweenOrEqual, 0, 1)
			}
		})

		Convey("new bits from bytes", func() {
			Convey("with full bytes", func() {
				bits := NewBitsFromBytes([]byte{0x00, 0xFF, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80})
				So(bits, ShouldResemble, Bits{
					0, 0, 0, 0, 0, 0, 0, 0,
					1, 1, 1, 1, 1, 1, 1, 1,
					1, 0, 0, 0, 0, 0, 0, 0,
					0, 1, 0, 0, 0, 0, 0, 0,
					0, 0, 1, 0, 0, 0, 0, 0,
					0, 0, 0, 1, 0, 0, 0, 0,
					0, 0, 0, 0, 1, 0, 0, 0,
					0, 0, 0, 0, 0, 1, 0, 0,
					0, 0, 0, 0, 0, 0, 1, 0,
					0, 0, 0, 0, 0, 0, 0, 1,
				})
			})

			Convey("with string", func() {
				str := NewBitsFromBytes([]byte("Hello!"))
				So(str, ShouldResemble, Bits{
					0, 0, 0, 1, 0, 0, 1, 0,
					1, 0, 1, 0, 0, 1, 1, 0,
					0, 0, 1, 1, 0, 1, 1, 0,
					0, 0, 1, 1, 0, 1, 1, 0,
					1, 1, 1, 1, 0, 1, 1, 0,
					1, 0, 0, 0, 0, 1, 0, 0,
				})
			})

		})

		Convey("to bytes", func() {
			Convey("with empty", func() {
				bits := Bits{}
				So(bits.ToBytes(), ShouldResemble, []byte{})
			})

			Convey("with 1 incomplete byte", func() {
				bits := Bits{1, 1, 1}
				So(bits.ToBytes(), ShouldResemble, []byte{0x07})
			})

			Convey("with 2 incomplete bytes", func() {
				bits := Bits{
					1, 1, 1, 1, 1, 1, 1, 1,
					1,
				}
				So(bits.ToBytes(), ShouldResemble, []byte{0xFF, 0x01})
			})

			Convey("with full bytes", func() {
				bits := Bits{
					0, 0, 0, 0, 0, 0, 0, 0,
					1, 1, 1, 1, 1, 1, 1, 1,
					1, 0, 0, 0, 0, 0, 0, 0,
					0, 1, 0, 0, 0, 0, 0, 0,
					0, 0, 1, 0, 0, 0, 0, 0,
					0, 0, 0, 1, 0, 0, 0, 0,
					0, 0, 0, 0, 1, 0, 0, 0,
					0, 0, 0, 0, 0, 1, 0, 0,
					0, 0, 0, 0, 0, 0, 1, 0,
					0, 0, 0, 0, 0, 0, 0, 1,
				}
				So(bits.ToBytes(), ShouldResemble, []byte{0x00, 0xFF, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80})
			})

			Convey("with string", func() {
				str := Bits{
					0, 0, 0, 1, 0, 0, 1, 0,
					1, 0, 1, 0, 0, 1, 1, 0,
					0, 0, 1, 1, 0, 1, 1, 0,
					0, 0, 1, 1, 0, 1, 1, 0,
					1, 1, 1, 1, 0, 1, 1, 0,
					1, 0, 0, 0, 0, 1, 0, 0,
				}
				So(str.ToBytes(), ShouldResemble, []byte("Hello!"))
			})
		})
	})
}
