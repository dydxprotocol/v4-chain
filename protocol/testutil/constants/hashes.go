package constants

// SHA256 hashes.
var (
	// Order hashes.
	OrderHash_Empty = [32]byte{
		0x06, 0xd9, 0x36, 0x73, 0x41, 0x3b, 0xa3, 0x29, 0x4a, 0xbc, 0x0e, 0x33, 0x97, 0x38, 0x1a, 0x44,
		0x15, 0xed, 0x37, 0x9c, 0xe6, 0x21, 0xfe, 0x91, 0x6b, 0xaf, 0xf7, 0x45, 0x98, 0x3d, 0x1b, 0x8c,
	}
	OrderHash_Alice_Number0_Id0 = [32]uint8{
		0xc5, 0xb5, 0xc0, 0x77, 0x3f, 0xfa, 0x4f, 0x7, 0x39, 0x61, 0x4f, 0x55, 0x87, 0xcd, 0x80, 0x96, 0xa3, 0x2d, 0x6c, 0xa7, 0x28, 0x53, 0xe7, 0x2b, 0x20, 0xad, 0x86, 0x75, 0xeb, 0x10, 0x46, 0x15,
	}

	// Liquidation order hashes.
	LiquidationOrderHash_Empty = [32]byte{
		0x10, 0x2b, 0x51, 0xb9, 0x76, 0x5a, 0x56, 0xa3, 0xe8, 0x99, 0xf7, 0xcf, 0x0e, 0xe3, 0x8e, 0x52,
		0x51, 0xf9, 0xc5, 0x03, 0xb3, 0x57, 0xb3, 0x30, 0xa4, 0x91, 0x83, 0xeb, 0x7b, 0x15, 0x56, 0x04,
	}
	LiquidationOrderHash_Alice_Number0_Perpetual0 = [32]byte{
		0x82, 0xba, 0x4e, 0xbd, 0x7b, 0x36, 0x58, 0x3c, 0x30, 0x37, 0x5b, 0x88, 0xb5, 0x9a, 0x8f, 0x34, 0x9b, 0x7a, 0x4e, 0xfa, 0x6e, 0xe3, 0x67, 0x65, 0x3c, 0xdf, 0x50, 0x40, 0x7e, 0xa3, 0x26, 0x27,
	}
	LiquidationOrderHash_Alice_Number0_Perpetual1 = [32]byte{
		0x92, 0x15, 0xab, 0xb0, 0xa5, 0xd5, 0x64, 0x79, 0x96, 0x39, 0xe9, 0x23, 0xcb, 0x1e, 0x67, 0x29, 0xbb, 0xc9, 0x56, 0x76, 0x82, 0xa4, 0x21, 0x8e, 0x58, 0xa6, 0xcd, 0xee, 0xbc, 0x3, 0xc2, 0xf9,
	}
	LiquidationOrderHash_Alice_Number1_Perpetual0 = [32]byte{
		0x97, 0xd7, 0x5a, 0x43, 0x53, 0x9c, 0x8d, 0x16, 0xcf, 0xa1, 0x24, 0xb9, 0x35, 0x2b, 0xa8, 0xdb, 0x95, 0x16, 0x8e, 0xf8, 0xd, 0xdd, 0x9d, 0x5e, 0x2f, 0xe, 0x42, 0x8c, 0xb7, 0x69, 0xdc, 0xfd,
	}

	// OrderId hashes
	OrderIdHash_Empty = []byte{
		0x10, 0x2b, 0x51, 0xb9, 0x76, 0x5a, 0x56, 0xa3, 0xe8, 0x99, 0xf7, 0xcf, 0x0e, 0xe3, 0x8e, 0x52,
		0x51, 0xf9, 0xc5, 0x03, 0xb3, 0x57, 0xb3, 0x30, 0xa4, 0x91, 0x83, 0xeb, 0x7b, 0x15, 0x56, 0x04,
	}
	OrderIdHash_Alice_Number0_Id0 = []byte{
		0x82, 0xba, 0x4e, 0xbd, 0x7b, 0x36, 0x58, 0x3c, 0x30, 0x37, 0x5b, 0x88, 0xb5, 0x9a, 0x8f, 0x34, 0x9b, 0x7a, 0x4e, 0xfa, 0x6e, 0xe3, 0x67, 0x65, 0x3c, 0xdf, 0x50, 0x40, 0x7e, 0xa3, 0x26, 0x27,
	}
)
