package main

import "eink-go-client/epd"

type Color uint16

func Bpp(bpp int) epd.PixelMode {
	switch bpp {
	case 2:
		return epd.BPP2
	case 3:
		return epd.BPP3
	case 4:
		return epd.BPP4
	case 1:
	case 8:
		return epd.BPP8
	}
	return epd.BPP8
}
