package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/pushover"
	"gopkg.in/yaml.v2"
)

const configFile = "config.yaml"

type Config struct {
	Server            string   `yaml:"server"`
	Topics            []string `yaml:"topics"`
	QoS               int      `yaml:"qos"`
	ClientID          string   `yaml:"clientid"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	PushoverToken     string   `yaml:"pushover_token"`
	PushoverRecipient string   `yaml:"pushover_recipient"`
}

var cachedConfig *Config

func readConfig() (*Config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	config := new(Config)

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func sendPushoverMessage(image []byte) {
	config, err := readConfig()
	if err != nil {
		log.Panic(err)
		//os.Exit(1)
	}
	app := pushover.New(config.PushoverToken)
	recipient := pushover.NewRecipient(config.PushoverRecipient)
	message := &pushover.Message{
		Message:  "Person detected",
		Title:    "Alarm",
		Priority: pushover.PriorityEmergency,
		//URL:         "http://google.com",
		//URLTitle:    "Google",
		Timestamp:  time.Now().Unix(),
		Retry:      60 * time.Second,
		Expire:     time.Hour,
		DeviceName: "Alienphone",
		//CallbackURL: "http://yourapp.com/callback",
		Sound: pushover.SoundSiren,
	}
	reader := bytes.NewReader(image)
	message.AddAttachment(reader)
	// Send the message to the recipient
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Panic(err)
		//os.Exit(1)
	}
	// Print the response if you want
	log.Println(response)
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	log.Printf("Received message on topic: %s", message.Topic())
	sendPushoverMessage(message.Payload())
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
		os.Exit(1)
	}
	//sendPushoverMessage(config)
	// Create an MQTT client options
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Server)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Could not connect")
		panic(token.Error())
	}

	// Subscribe to each topic individually
	for _, topic := range config.Topics {
		if token := client.Subscribe(topic, byte(config.QoS), onMessageReceived); token.Wait() && token.Error() != nil {
			log.Printf("Error subscribing to topic %s: %v", topic, token.Error())
		}
	}

	// Check if the client successfully connected to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
		return
	}

	log.Println("Connected to MQTT broker. Subscribed to topics:", strings.Join(config.Topics, ", "))

	// Wait for a termination signal to gracefully disconnect
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<-sigChannel

	// Unsubscribe and disconnect from the MQTT broker
	for _, topic := range config.Topics {
		client.Unsubscribe(topic)
	}
	client.Disconnect(250)
}
