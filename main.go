package main

import (
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DeviceConfig represents device configuration.
type DeviceConfig struct {
	// Username aka device_id
	Username string

	// Password aka primaryKey
	Password string

	// JSONStPath represents JSON mc_st file path.
	JSONStPath string

	// JSONOpPath represents JSON mc_op file path.
	JSONOpPath string

	// JSONRopPath represents JSON mc_rop file path.
	JSONRopPath string

	// Topic represents MQTT topic.
	Topic string
}

// Message represents message schema.
type Message struct {
	Name  string    `json:"name"`
	Value float64   `json:"value"`
	Time  time.Time `json:"time"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Info().
				Err(err).
				Msg(".env file not found. Use shell environment variables instead")
		} else {
			log.Fatal().
				Err(err).
				Msg("Load .env file failed")
		}
	}

	brokerURL := os.Getenv("MDG_BROKER_URL")

	byteDevConfs, err := ioutil.ReadFile("device-configs.json")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Load device-configs.json file failed")
	}

	var devConfs []DeviceConfig

	err = json.Unmarshal(byteDevConfs, &devConfs)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unmarshal failed")
	}

	for i := 0; i < len(devConfs); i++ {
		conf := devConfs[i]

		go func() {
			opts := mqtt.NewClientOptions().
				AddBroker(brokerURL).
				SetUsername(conf.Username).
				SetPassword(conf.Password).
				SetClientID(conf.Username)
			mqttClient := mqtt.NewClient(opts)

			log.Info().
				Str("url", brokerURL).
				Str("username", conf.Username).
				Msg("Connect to MQTT broker")

			if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
				log.Fatal().
					Err(token.Error()).
					Str("url", brokerURL).
					Str("username", conf.Username).
					Msg("Connect to MQTT broker failed")
			}

			byteStMsgs, err := ioutil.ReadFile(conf.JSONStPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read status messages file failed")
			}

			var stMsgs []Message

			err = json.Unmarshal(byteStMsgs, &stMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal status messages failed")
			}

			byteOpMsgs, err := ioutil.ReadFile(conf.JSONOpPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read output messages file failed")
			}

			var opMsgs []Message

			err = json.Unmarshal(byteOpMsgs, &opMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal output messages failed")
			}

			byteRopMsgs, err := ioutil.ReadFile(conf.JSONOpPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read reject output messages file failed")
			}

			var ropMsgs []Message

			err = json.Unmarshal(byteRopMsgs, &ropMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal reject output messages failed")
			}

			token := mqttClient.Publish(conf.Topic, 2, false, nil)
			if token.Wait() && token.Error() != nil {
				log.Error().
					Err(token.Error()).
					Msg("Publish MQTT message failed")
			}

			// Make the program keep running until receive an interruption signal.
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			// Abort the program upon receiving an interruption signal.
			s := <-c
			log.Info().Msgf("Received %s signal. Aborting", s)

			mqttClient.Disconnect(250)
		}()
	}
}
