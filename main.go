package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Publisher represents publisher configuration.
type Publisher struct {
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

// IfraMessage represents Ifra MQTT broker compatible message.
type IfraMessage struct {
	N string  `json:"n"`
	V float64 `json:"v"`
	U string  `json:"u"`
}

func main() {
	// Customize usage.
	flag.Usage = func() {
		//nolint:forbidigo // Print usage.
		fmt.Printf("Usage: machine-data-generator")
		flag.PrintDefaults()
	}

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Info().
				Err(err).
				Msg(".env file not found. Use shell environment variables")
		} else {
			log.Fatal().
				Err(err).
				Msg("Read .env file failed")
		}
	}

	brokerURL := os.Getenv("MDG_BROKER_URL")

	bytePublConfs, err := ioutil.ReadFile("data/publishers.json")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Read publishers.json file failed")
	}

	var publishers []Publisher

	err = json.Unmarshal(bytePublConfs, &publishers)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unmarshal failed")
	}

	for i := 0; i < len(publishers); i++ {
		publisher := publishers[i]

		var (
			qos      byte = 2
			retained bool
		)

		go func() {
			opts := mqtt.NewClientOptions().
				AddBroker(brokerURL).
				SetUsername(publisher.Username).
				SetPassword(publisher.Password).
				SetClientID(publisher.Username)
			mqttClient := mqtt.NewClient(opts)

			log.Info().
				Str("url", brokerURL).
				Str("username", publisher.Username).
				Msg("Connect to MQTT broker")

			if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
				log.Fatal().
					Err(token.Error()).
					Str("url", brokerURL).
					Str("username", publisher.Username).
					Msg("Connect to MQTT broker failed")
			}

			byteStMsgs, err := ioutil.ReadFile(publisher.JSONStPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read status messages file failed")
			}

			var stMsgs []IfraMessage

			err = json.Unmarshal(byteStMsgs, &stMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal status messages failed")
			}

			byteOpMsgs, err := ioutil.ReadFile(publisher.JSONOpPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read output messages file failed")
			}

			var opMsgs []IfraMessage

			err = json.Unmarshal(byteOpMsgs, &opMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal output messages failed")
			}

			byteRopMsgs, err := ioutil.ReadFile(publisher.JSONRopPath)
			if err != nil {
				log.Fatal().Err(err).Msg("Read reject output messages file failed")
			}

			var ropMsgs []IfraMessage

			err = json.Unmarshal(byteRopMsgs, &ropMsgs)
			if err != nil {
				log.Fatal().Err(err).Msg("Unmarshal reject output messages failed")
			}

			// Publish status message.
			go func() {
			Start:
				for _, msg := range stMsgs {
					ifraMsg := []IfraMessage{msg}
					m, _ := json.Marshal(ifraMsg)

					token := mqttClient.Publish(publisher.Topic, qos, retained, m)
					if token.Wait() && token.Error() != nil {
						log.Error().Err(token.Error()).Msg("Publish status message failed")
					}

					log.Info().
						Str("to", brokerURL).
						Str("topic", publisher.Topic).
						Uint8("qos", qos).
						Bool("retained", retained).
						RawJSON("msg", m).
						Msg("Publish MQTT message")

					time.Sleep(time.Second)
				}

				goto Start
			}()

			// Publish output message.
			go func() {
			Start:
				for _, msg := range opMsgs {
					ifraMsg := []IfraMessage{msg}
					m, _ := json.Marshal(ifraMsg)

					token := mqttClient.Publish(publisher.Topic, qos, retained, m)
					if token.Wait() && token.Error() != nil {
						log.Error().Err(token.Error()).Msg("Publish output message failed")
					}

					log.Info().
						Str("to", brokerURL).
						Str("topic", publisher.Topic).
						Uint8("qos", qos).
						Bool("retained", retained).
						RawJSON("msg", m).
						Msg("Publish MQTT message")

					time.Sleep(time.Second)
				}

				goto Start
			}()

			// Publish reject output message.
			go func() {
			Start:
				for _, msg := range ropMsgs {
					ifraMsg := []IfraMessage{msg}
					m, _ := json.Marshal(ifraMsg)

					token := mqttClient.Publish(publisher.Topic, qos, retained, m)
					if token.Wait() && token.Error() != nil {
						log.Error().Err(token.Error()).Msg("Publish reject output message failed")
					}

					log.Info().
						Str("to", brokerURL).
						Str("topic", publisher.Topic).
						Uint8("qos", qos).
						Bool("retained", retained).
						RawJSON("msg", m).
						Msg("Publish MQTT message")

					time.Sleep(time.Second)
				}

				goto Start
			}()

			// Make the program keep running until receive an interruption signal.
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			// Abort the program upon receiving an interruption signal.
			<-c

			mqttClient.Disconnect(250)
		}()
	}

	// Make the program keep running until receive an interruption signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Abort the program upon receiving an interruption signal.
	s := <-c
	log.Info().Msgf("Received %s signal. Aborting", s)
}
