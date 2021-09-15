package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	brokerURL := "tcp://staging.mqtt.ifra.io:1883"
	//topic := "organization/93d818ed-ea26-4ccf-a9ef-e1852197e96f/messages"
	username := "ecbe49d1-8b2d-4f16-afb4-3e323a74d214"
	password := "5b76e53f-24a1-485b-8284-d55eef426234"

	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetUsername(username).
		SetPassword(password)
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().
			Err(token.Error()).
			Str("url", brokerURL).
			Msg("Connect to source MQTT broker failed")
	}
}
