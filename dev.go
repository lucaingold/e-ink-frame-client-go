package main

import (
	"log"
	"time"
)

//

const (
	EpdRstPin  = 17 //11 // Raspberry Pi Pin 17
	EpdCsPin   = 8  //24 // Raspberry Pi Pin 8
	EpdBusyPin = 24 //18 // Raspberry Pi Pin 24
)

var (
	rstPin   Pin
	csPin    Pin
	readyPin Pin
)

// Open sets the I/O ports and SPI
func Open() (err error) {
	Debug("Init start")

	if err := OpenRpio(); err != nil {
		log.Fatalln("RPIO Open Error:", err)
	}

	//
	// init SPI
	//

	Debug("Initializing SPI")

	if err := SpiBegin(Spi0); err != nil {
		log.Fatalln("SpiBegin Error:", err)
	}

	SpiChipSelect(0)
	SpiSpeed(24000000) // 24MHz
	SpiMode(0, 0)

	//
	// init pins
	//

	Debug("Initializing GPIO pins")

	rstPin = Pin(EpdRstPin)
	csPin = Pin(EpdCsPin)
	readyPin = Pin(EpdBusyPin)

	rstPin.Output()
	csPin.Output()
	readyPin.Input()

	csOff()

	Debug("EPD initialization complete")
	return nil
}

// Close ends SPI usage and restores pins
func Close() {
	Debug("Shutting down EPD")
	csPin.Low()
	rstPin.Low()

	SpiEnd(Spi0)

	CloseRpio()
}

// csOn selects slave
func csOn() {
	//Debug("CS On")
	csPin.Low()
}

// csOff deselects slave
func csOff() {
	//Debug("CS Off")
	csPin.High()
}

// Reset resets a slave
func Reset() {
	Debug("EPD Reset")
	rstPin.High()
	time.Sleep(time.Duration(200) * time.Millisecond)
	rstPin.Low()
	time.Sleep(time.Duration(10) * time.Millisecond)
	rstPin.High()
	time.Sleep(time.Duration(200) * time.Millisecond)
}
