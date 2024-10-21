package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var devInfo *DevInfo // Global variable for device information

func main() {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// Read mqtt configuration values
	mqttConfig := viper.GetStringMapString("mqtt")

	topic := viper.GetString("mqtt.topic")

	c, err := newClient(mqttConfig)

	if err != nil {
		panic(fmt.Errorf("failed to create MQTT client: %v", err))
	}

	if c == nil {
		panic("MQTT client is nil")
	}

	vcom := viper.GetUint16("display.vcom")
	devInfo := Init(vcom)
	fmt.Println(devInfo)
	defer Exit()

	if err := c.Subscribe(topic, func(_ MQTT.Client, m MQTT.Message) {
		fmt.Printf("Message: %s \n", m.Payload())
		fmt.Printf("Topic: %s \n", m.Topic())

		fmt.Printf("Received message on topic %s\n", m.Topic())
		imageBytes := m.Payload() // Get the byte array from the message

		// Create a DataBuffer from the received byte array
		var imageBuffer DataBuffer
		imageBuffer = createDataBuffer(imageBytes) // Implement this function to convert byte array to DataBuffer

		// Call the displayImage function
		displayImage(imageBuffer, 0, 0, devInfo.PanelW, devInfo.PanelH)
	}); err != nil {
		panic(err)
	}

	if err := c.Publish("Hello World", topic); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	c.mqttClient.Unsubscribe(topic)
	c.mqttClient.Disconnect(250)
}

func displayImage(imageBuffer DataBuffer, x, y, width, height uint16) {
	// Set up load image info
	imageInfo := LoadImgInfo{
		EndianType:       LoadImgLittleEndian,
		PixelFormat:      BPP4,                    // Assuming 4 bits per pixel
		Rotate:           Rotate0,                 // No rotation
		SourceBufferAddr: imageBuffer,             // Your image data
		TargetMemAddr:    devInfo.TargetAddress(), // Target memory address from device info
	}

	// Define the area to display
	areaInfo := AreaImgInfo{
		X: x,
		Y: y,
		W: width,
		H: height,
	}

	// Load image and display it
	imageInfo.HostAreaPackedPixelWrite(areaInfo, 4, true)
	DisplayArea(x, y, width, height, GC16Mode)
}

func createDataBuffer(imageBytes []byte) DataBuffer {
	// Ensure the byte slice can be converted to uint16
	if len(imageBytes)%2 != 0 {
		fmt.Println("Warning: Image byte length is not even, truncating last byte.")
		imageBytes = imageBytes[:len(imageBytes)-1] // Truncate if odd
	}

	// Create a DataBuffer of the appropriate size
	buffer := make(DataBuffer, len(imageBytes)/2)

	// Convert byte array to uint16
	for i := 0; i < len(imageBytes); i += 2 {
		buffer[i/2] = uint16(imageBytes[i]) | (uint16(imageBytes[i+1]) << 8)
	}

	return buffer
}