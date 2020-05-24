package oled

var font [][]byte
var signalLevels [][]byte

func init() {
	font = [][]byte{
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00}, // Space
		[]byte{0x00, 0x00, 0xBE, 0x00, 0x00}, // !
		[]byte{0x00, 0x06, 0x00, 0x06, 0x00}, // "
		[]byte{0x28, 0xFE, 0x28, 0xFE, 0x28}, // #
		[]byte{0x00, 0x5C, 0xFE, 0x74, 0x00}, // $
		[]byte{0x0C, 0x2C, 0x10, 0x68, 0x60}, // %
		[]byte{0x48, 0xB4, 0xA4, 0x40, 0xA0}, // &
		[]byte{0x00, 0x00, 0x06, 0x00, 0x00}, // '
		[]byte{0x00, 0x00, 0x7C, 0x82, 0x00}, // (
		[]byte{0x00, 0x82, 0x7C, 0x00, 0x00}, // )
		[]byte{0x6C, 0x38, 0x7C, 0x38, 0x6C}, // *
		[]byte{0x10, 0x10, 0x7C, 0x10, 0x10}, // +
		[]byte{0x00, 0x80, 0xE0, 0x60, 0x00}, // ,
		[]byte{0x10, 0x10, 0x10, 0x10, 0x10}, // -
		[]byte{0x00, 0x00, 0xC0, 0xC0, 0x00}, // .
		[]byte{0x80, 0x60, 0x10, 0x0C, 0x02}, // /
		[]byte{0x7C, 0xA2, 0x92, 0x8A, 0x7C}, // 0
		[]byte{0x00, 0x84, 0xFE, 0x80, 0x00}, // 1
		[]byte{0xE4, 0x92, 0x92, 0x92, 0x8C}, // 2
		[]byte{0x44, 0x92, 0x92, 0x92, 0x6C}, // 3
		[]byte{0x30, 0x28, 0x24, 0xFE, 0x20}, // 4
		[]byte{0x5E, 0x92, 0x92, 0x92, 0x62}, // 5
		[]byte{0x7C, 0x92, 0x92, 0x92, 0x64}, // 6
		[]byte{0x06, 0x02, 0xE2, 0x12, 0x0E}, // 7
		[]byte{0x6C, 0x92, 0x92, 0x92, 0x6C}, // 8
		[]byte{0x4C, 0x92, 0x92, 0x92, 0x7C}, // 9
		[]byte{0x00, 0x00, 0x6C, 0x6C, 0x00}, // :
		[]byte{0x00, 0x80, 0xEC, 0x6C, 0x00}, // ;
		[]byte{0x10, 0x28, 0x28, 0x44, 0x44}, // <
		[]byte{0x28, 0x28, 0x28, 0x28, 0x28}, // =
		[]byte{0x44, 0x44, 0x28, 0x28, 0x10}, // >
		[]byte{0x04, 0x02, 0xA2, 0x12, 0x0C}, // ?
		[]byte{0x7C, 0x82, 0xBA, 0xAA, 0xBC}, // @
		[]byte{0xF8, 0x24, 0x22, 0x24, 0xF8}, // A
		[]byte{0xFE, 0x92, 0x92, 0x92, 0x6C}, // B
		[]byte{0x7C, 0x82, 0x82, 0x82, 0x44}, // C
		[]byte{0x82, 0xFE, 0x82, 0x82, 0x7C}, // D
		[]byte{0xFE, 0x92, 0x92, 0x92, 0x82}, // E
		[]byte{0xFE, 0x12, 0x12, 0x12, 0x02}, // F
		[]byte{0x7C, 0x82, 0x82, 0xA2, 0x64}, // G
		[]byte{0xFE, 0x10, 0x10, 0x10, 0xFE}, // H
		[]byte{0x00, 0x82, 0xFE, 0x82, 0x00}, // I
		[]byte{0x40, 0x80, 0x82, 0x7E, 0x02}, // J
		[]byte{0xFE, 0x10, 0x28, 0x44, 0x82}, // K
		[]byte{0xFE, 0x80, 0x80, 0x80, 0x80}, // L
		[]byte{0xFE, 0x0C, 0x18, 0x0C, 0xFE}, // M
		[]byte{0xFE, 0x0C, 0x38, 0x60, 0xFE}, // N
		[]byte{0x7C, 0x82, 0x82, 0x82, 0x7C}, // O
		[]byte{0xFE, 0x12, 0x12, 0x12, 0x0C}, // P
		[]byte{0x7C, 0x82, 0xA2, 0xC2, 0xFC}, // Q
		[]byte{0xFE, 0x12, 0x32, 0x52, 0x8C}, // R
		[]byte{0x4C, 0x92, 0x92, 0x92, 0x64}, // S
		[]byte{0x02, 0x02, 0xFE, 0x02, 0x02}, // T
		[]byte{0x7E, 0x80, 0x80, 0x80, 0x7E}, // U
		[]byte{0x1E, 0x60, 0x80, 0x60, 0x1E}, // V
		[]byte{0x7E, 0x80, 0x70, 0x80, 0x7E}, // W
		[]byte{0xC6, 0x28, 0x10, 0x28, 0xC6}, // X
		[]byte{0x0E, 0x10, 0xE0, 0x10, 0x0E}, // Y
		[]byte{0xC6, 0xA2, 0x92, 0x8A, 0xC6}, // Z
		[]byte{0x00, 0x00, 0xFE, 0x82, 0x00}, // [
		[]byte{0x02, 0x0C, 0x10, 0x60, 0x80}, // \
		[]byte{0x00, 0x82, 0xFE, 0x00, 0x00}, // ]
		[]byte{0x1C, 0x3E, 0xFC, 0x3E, 0x1C}, // ♥ (^)
		[]byte{0x80, 0x80, 0x80, 0x80, 0x80}, // _
		[]byte{0x00, 0x02, 0x06, 0x00, 0x00}, // `
		[]byte{0xFE, 0x82, 0x92, 0x82, 0xFE}, // Unknown
	}
	signalLevels = [][]byte{
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x00, 0x00, 0x40, 0x20, 0xA0, 0x20, 0x40, 0x00, 0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x20, 0x10, 0x48, 0x28, 0xA4, 0x28, 0x48, 0x10, 0x20, 0x00, 0x00},
		[]byte{0x10, 0x08, 0x24, 0x12, 0x4A, 0x29, 0xA5, 0x29, 0x4A, 0x12, 0x24, 0x08, 0x10},
	}
	SignalLevels = len(signalLevels)
}