package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type client struct {
	mqttClient MQTT.Client
}

func newClient(mqttConfig map[string]string) (*client, error) {
	broker := mqttConfig["broker"]
	clientID := mqttConfig["client_id"]
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

	return &client{
		mqttClient,
	}, nil
}

func (c *client) Publish(msg, topic string) error {
	if token := c.mqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *client) Subscribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
