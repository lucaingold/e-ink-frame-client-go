package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type Client struct {
	MqttClient MQTT.Client
}

func NewClient(mqttConfig map[string]string) (*Client, error) {
	fmt.Printf("Creating MQTT client with ID: %s\n", mqttConfig["client_id"])
	broker := mqttConfig["broker"]
	clientID := uuid.New().String()
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%s", broker, mqttConfig["port"]))
	opts.SetClientID(clientID)
	opts.SetUsername(mqttConfig["username"])
	opts.SetPassword(mqttConfig["password"])
	//opts.SetTLSConfig(&tls.Config{
	//	InsecureSkipVerify: true,
	//})

	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{
		mqttClient,
	}, nil
}

func (c *Client) Publish(msg, topic string) error {
	if token := c.MqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Subscribe(topic string, f MQTT.MessageHandler) error {
	if token := c.MqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
