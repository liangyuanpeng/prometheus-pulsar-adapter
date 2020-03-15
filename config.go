package main

import (
	"os"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/sirupsen/logrus"
)

var (
	pulsarBroker = "pulsar://pulsar:6650"
	pulsarTopic  = "metrics"
	serializer   Serializer
	client       pulsar.Client
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if value := os.Getenv("LOG_LEVEL"); value != "" {
		logrus.SetLevel(parseLogLevel(value))
	}

	if value := os.Getenv("PULSAR_BROKER"); value != "" {
		pulsarBroker = value
	}

	if value := os.Getenv("PULSAR_TOPIC"); value != "" {
		pulsarTopic = value
	}

	var authentication pulsar.Authentication

	if value := os.Getenv("PULSAR_AUTH_TOKEN"); value != "" {
		authentication = pulsar.NewAuthenticationToken(value)
	}

	var err error
	serializer, err = NewJSONSerializer()
	if err != nil {
		logrus.WithError(err).Error("couldn't init json")
	}

	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:            pulsarBroker,
		Authentication: authentication,
	})
	if err != nil {
		logrus.WithError(err).Error("couldn't init json")
	}
}

func parseLogLevel(value string) logrus.Level {
	level, err := logrus.ParseLevel(value)

	if err != nil {
		logrus.WithField("log-level-value", value).Warningln("invalid log level from env var, using info")
		return logrus.InfoLevel
	}

	return level
}
