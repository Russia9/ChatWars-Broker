package main

import (
	"cw-broker/lib"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger = logrus.New()

func main() {
	// Logger init
	logger.Out = os.Stdout
	logger.Info("Initializing ChatWars Broker")

	// Change logger log level
	switch os.Getenv("CWBR_LOGLEVEL") {
	case "TRACE":
		logger.SetLevel(logrus.TraceLevel)
		break
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
		break
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
		break
	case "WARN":
		logger.SetLevel(logrus.WarnLevel)
		break
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
		break
	case "FATAL":
		logger.SetLevel(logrus.FatalLevel)
		break
	case "PANIC":
		logger.SetLevel(logrus.PanicLevel)
		break
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	SentryDSN := lib.GetEnv("CWBR_SENTRY_DSN", "")
	SentryEnvironment := lib.GetEnv("CWBR_ENVIRONMENT", "production")

	// Sentry init
	logger.Debug("Initializing Sentry")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         SentryDSN,
		Environment: SentryEnvironment,
	})
	if err != nil {
		logger.Warn("Sentry init error: ", err.Error())
	}
	defer sentry.Flush(2 * time.Second)

	// Kafka consumer init
	logger.Debug("Initializing Kafka")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": lib.GetEnv("CWBR_KAFKA_ADDRESS", "localhost"),
		"group.id":          "cw3",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		logger.Panic(err)
		sentry.CaptureException(err)
		return
	}

	err = InitBot(os.Getenv("CWBR_BOT_TOKEN"), logger, consumer)
	if err != nil {
		sentry.CaptureException(err)
		logger.Panic(err)
	}
}
