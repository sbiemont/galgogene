package gene

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerializer(t *testing.T) {
	Convey("to bits", t, func() {
		Convey("when n=8", func() {
			bytes := []byte{0x00, 0xFF, 0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
			szr := Serializer{
				pack: 8,
			}

			So(szr.ToBits(bytes), ShouldResemble, Bits{
				MaxValue: 1,
				Raw: []uint8{
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
				},
			})
		})

		Convey("when n=3", func() {
			bytes := []byte{0b00000000, 0b00000001, 0b00000010, 0b00000011}
			szr := Serializer{
				pack: 3,
			}

			So(szr.ToBits(bytes), ShouldResemble, Bits{
				MaxValue: 1,
				Raw: []uint8{
					0, 0, 0,
					0, 0, 1,
					0, 1, 0,
					0, 1, 1,
				},
			})
		})

		Convey("when n=1", func() {
			// Only the last bit is gathered
			bytes := []byte{0x00, 0x01, 0xFF, 0xFE}
			szr := Serializer{
				pack: 1,
			}

			So(szr.ToBits(bytes), ShouldResemble, Bits{
				MaxValue: 1,
				Raw: []uint8{
					0,
					1,
					1,
					0,
				},
			})
		})
	})

	Convey("to bytes", t, func() {
		Convey("when error nb bits", func() {
			bits := newBits([]uint8{0, 0, 0})
			szr, err := NewSerializer(2)
			So(err, ShouldBeNil)

			by, err := szr.ToBytes(bits)
			So(err, ShouldBeError, "cannot group, total nb of bits (3) should be modulo 2")
			So(by, ShouldBeNil)
		})

		Convey("when error n>8", func() {
			szr, err := NewSerializer(16)
			So(err, ShouldBeError, "cannot group, n > 8")
			So(szr, ShouldResemble, Serializer{})
		})

		Convey("when ok", func() {
			bits := newBits([]uint8{0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0})

			Convey("group by 2", func() {
				szr, err := NewSerializer(2)
				So(err, ShouldBeNil)

				by, err := szr.ToBytes(bits)
				So(err, ShouldBeNil)
				So(by, ShouldResemble, []uint8{0b00, 0b01, 0b10, 0b11, 0b11, 0b10, 0b01, 0b00})
			})

			Convey("group by 4", func() {
				szr, err := NewSerializer(4)
				So(err, ShouldBeNil)

				by, err := szr.ToBytes(bits)
				So(err, ShouldBeNil)
				So(by, ShouldResemble, []uint8{0b0001, 0b1011, 0b1110, 0b0100})
			})

			Convey("group by 8", func() {
				szr, err := NewSerializer(8)
				So(err, ShouldBeNil)

				by, err := szr.ToBytes(bits)
				So(err, ShouldBeNil)
				So(by, ShouldResemble, []uint8{0b00011011, 0b11100100})
			})
		})
	})
}
