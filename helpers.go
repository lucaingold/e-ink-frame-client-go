package main

type Color uint16

func Bpp(bpp int) PixelMode {
	switch bpp {
	case 2:
		return BPP2
	case 3:
		return BPP3
	case 4:
		return BPP4
	case 1:
	case 8:
		return BPP8
	}
	return BPP8
}
