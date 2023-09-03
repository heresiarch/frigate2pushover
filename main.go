package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/pushover"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server            string   `yaml:"server"`
	Topic             []string `yaml:"topic"`
	QoS               int      `yaml:"qos"`
	ClientID          string   `yaml:"clientid"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	PushoverToken     string   `yaml:"pushover_token"`
	PushoverRecipient string   `yaml:"pushover_recipient"`
}

func readConfig(configFile string) (Config, error) {
	var config Config

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

func sendPushoverMessage(config Config) {
	log.Printf("%+v\n", config)
	app := pushover.New(config.PushoverToken)

	recipient := pushover.NewRecipient(config.PushoverRecipient)

	message := &pushover.Message{
		Message:     "My awesome message",
		Title:       "My title",
		Priority:    pushover.PriorityEmergency,
		URL:         "http://google.com",
		URLTitle:    "Google",
		Timestamp:   time.Now().Unix(),
		Retry:       60 * time.Second,
		Expire:      time.Hour,
		DeviceName:  "Alienphone",
		CallbackURL: "http://yourapp.com/callback",
		Sound:       pushover.SoundCosmic,
	}

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
	log.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
}

func main() {
	configFile := "config.yaml"
	config, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
		os.Exit(1)
	}
	//sendPushoverMessage(config)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	connOpts := MQTT.NewClientOptions().AddBroker(config.Server).SetClientID(*&config.ClientID).SetCleanSession(true)
	if config.Username != "" {
		connOpts.SetUsername(config.Username)
		if config.Password != "" {
			connOpts.SetPassword(config.Password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(config.Topic[0], byte(config.QoS), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("Connected to %s\n", config.Server)
	}
	<-c
}
