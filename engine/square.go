package engine

type Square uint8

var (
	A1      Square = 0x00
	B1      Square = 0x01
	C1      Square = 0x02
	D1      Square = 0x03
	E1      Square = 0x04
	F1      Square = 0x05
	G1      Square = 0x06
	H1      Square = 0x07
	A2      Square = 0x10
	B2      Square = 0x11
	C2      Square = 0x12
	D2      Square = 0x13
	E2      Square = 0x14
	F2      Square = 0x15
	G2      Square = 0x16
	H2      Square = 0x17
	A3      Square = 0x20
	B3      Square = 0x21
	C3      Square = 0x22
	D3      Square = 0x23
	E3      Square = 0x24
	F3      Square = 0x25
	G3      Square = 0x26
	H3      Square = 0x27
	A4      Square = 0x30
	B4      Square = 0x31
	C4      Square = 0x32
	D4      Square = 0x33
	E4      Square = 0x34
	F4      Square = 0x35
	G4      Square = 0x36
	H4      Square = 0x37
	A5      Square = 0x40
	B5      Square = 0x41
	C5      Square = 0x42
	D5      Square = 0x43
	E5      Square = 0x44
	F5      Square = 0x45
	G5      Square = 0x46
	H5      Square = 0x47
	A6      Square = 0x50
	B6      Square = 0x51
	C6      Square = 0x52
	D6      Square = 0x53
	E6      Square = 0x54
	F6      Square = 0x55
	G6      Square = 0x56
	H6      Square = 0x57
	A7      Square = 0x60
	B7      Square = 0x61
	C7      Square = 0x62
	D7      Square = 0x63
	E7      Square = 0x64
	F7      Square = 0x65
	G7      Square = 0x66
	H7      Square = 0x67
	A8      Square = 0x70
	B8      Square = 0x71
	C8      Square = 0x72
	D8      Square = 0x73
	E8      Square = 0x74
	F8      Square = 0x75
	G8      Square = 0x76
	H8      Square = 0x77
	Invalid Square = 0x80
)

var (
	SquareMap = map[Square]string{
		0x00: "a1", 0x01: "b1", 0x02: "c1", 0x03: "d1", 0x04: "e1", 0x05: "f1", 0x06: "g1", 0x07: "h1",
		0x10: "a2", 0x11: "b2", 0x12: "c2", 0x13: "d2", 0x14: "e2", 0x15: "f2", 0x16: "g2", 0x17: "h2",
		0x20: "a3", 0x21: "b3", 0x22: "c3", 0x23: "d3", 0x24: "e3", 0x25: "f3", 0x26: "g3", 0x27: "h3",
		0x30: "a4", 0x31: "b4", 0x32: "c4", 0x33: "d4", 0x34: "e4", 0x35: "f4", 0x36: "g4", 0x37: "h4",
		0x40: "a5", 0x41: "b5", 0x42: "c5", 0x43: "d5", 0x44: "e5", 0x45: "f5", 0x46: "g5", 0x47: "h5",
		0x50: "a6", 0x51: "b6", 0x52: "c6", 0x53: "d6", 0x54: "e6", 0x55: "f6", 0x56: "g6", 0x57: "h6",
		0x60: "a7", 0x61: "b7", 0x62: "c7", 0x63: "d7", 0x64: "e7", 0x65: "f7", 0x66: "g7", 0x67: "h7",
		0x70: "a8", 0x71: "b8", 0x72: "c8", 0x73: "d8", 0x74: "e8", 0x75: "f8", 0x76: "g8", 0x77: "h8",
	}

	SquareLookup = map[string]Square{
		"a1": A1, "a2": A2, "a3": A3, "a4": A4, "a5": A5, "a6": A6, "a7": A7, "a8": A8,
		"b1": B1, "b2": B2, "b3": B3, "b4": B4, "b5": B5, "b6": B6, "b7": B7, "b8": B8,
		"c1": C1, "c2": C2, "c3": C3, "c4": C4, "c5": C5, "c6": C6, "c7": C7, "c8": C8,
		"d1": D1, "d2": D2, "d3": D3, "d4": D4, "d5": D5, "d6": D6, "d7": D7, "d8": D8,
		"e1": E1, "e2": E2, "e3": E3, "e4": E4, "e5": E5, "e6": E6, "e7": E7, "e8": E8,
		"f1": F1, "f2": F2, "f3": F3, "f4": F4, "f5": F5, "f6": F6, "f7": F7, "f8": F8,
		"g1": G1, "g2": G2, "g3": G3, "g4": G4, "g5": G5, "g6": G6, "g7": G7, "g8": G8,
		"h1": H1, "h2": H2, "h3": H3, "h4": H4, "h5": H5, "h6": H6, "h7": H7, "h8": H8,
	}

	board64square = []uint8{
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	}
)

func square(rank, file int8) int8 {
	return (rank << 4) | file
}

func rank(square int8) int8 {
	return square >> 4
}

func file(square int8) int8 {
	return square & 7
}
